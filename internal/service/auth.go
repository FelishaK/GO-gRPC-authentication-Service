package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/FelishaK/grpcAuth/internal/config"
	"github.com/FelishaK/grpcAuth/internal/domain"
	authdto "github.com/FelishaK/grpcAuth/internal/dto/auth"
	"github.com/FelishaK/grpcAuth/internal/lib/jwt"
	repo "github.com/FelishaK/grpcAuth/internal/repository"
	jwtvalidate "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserWasNotFound = errors.New("user wasn't found")
	ErrUserExists = errors.New("user with this email already exists")
	ErrAddUser = errors.New("error while creating a user")
	ErrCreatingToken = errors.New("error while creating a token")
	ErrDecodingToken = errors.New("error while decoding a token")
	ErrValidatingToken = errors.New("error while validating a token")
	ErrTokenExpired = errors.New("token has been expired")
)

type AuthService interface {
	Register(context.Context, *authdto.RegisterRequest) (string, error)
	Login(context.Context, *authdto.LoginRequest) (*authdto.LoginResponse, error)
	GetAuthedUser(context.Context, *authdto.Token) (*authdto.AuthedUserResponse, error)
	RefreshToken(context.Context, *authdto.Token) (*authdto.Token, error)
}

type DefaultAuthService struct {
	cfg *config.Config
	log *slog.Logger
	repo repo.AuthRepository
}

func NewAuthService(log *slog.Logger, cfg *config.Config, repo repo.AuthRepository) *DefaultAuthService {
	return &DefaultAuthService{
		cfg: cfg,
		log: log,
		repo: repo,
	}
}

func (as *DefaultAuthService) Register(ctx context.Context, rq *authdto.RegisterRequest) (string, error) {
	const op = "authService.Register"

	//validating requestDTO
	err := rq.Validate()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rq.Password), 10)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	defaultRole := "user"

	//creating domain object
	user := domain.User{
		FirstName: rq.FirstName,
		LastName: rq.LastName,
		PasswordHash: passwordHash,
		Email: rq.Email,
		RegistrationDate: primitive.NewDateTimeFromTime(time.Now()),
		Gender: rq.Gender,
		Role: &defaultRole,
	}

	//insert a user if it doesn't exists yet
	if _, err := as.repo.FindUserByEmail(ctx, rq.Email); errors.Is(err, repo.ErrNoUserFound) {
		userId, err := as.repo.AddUser(ctx, user)

		if err != nil || userId == "" {
			return "", fmt.Errorf("%s: %w", op, ErrAddUser)
		}
		return userId, nil
	} else if errors.Is(err, repo.ErrUserExists) {
		return "", fmt.Errorf("%s: %w", op, ErrUserExists)
	}
	return "", fmt.Errorf("%s: %w", op, err)
}

func (as *DefaultAuthService) Login(ctx context.Context, lr *authdto.LoginRequest) (*authdto.LoginResponse, error) {
	const op = "service.Login"

	err := lr.Validate()
	if err != nil {
		as.log.Debug("DTO conversion error", slog.String(op, err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	user, err := as.repo.FindUserByEmail(ctx, lr.Email)

	if err != nil {
		as.log.Debug("User query error", slog.String(op, err.Error()))
		if errors.Is(err, repo.ErrNoUserFound) {
			as.log.Debug("User was not found", slog.String(op, err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserWasNotFound)
		}
		as.log.Debug("Error while querying for User", slog.String(op, err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err) 
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(lr.Password)); err != nil {
		as.log.Debug("Passwords do not match", slog.String(op, err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	aDuration, err := time.ParseDuration(as.cfg.AccessTokenTTL)
	if err != nil {
		as.log.Error("Invalid duration format", slog.String("op", op))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	accessToken, _ := jwt.NewToken(user.Id.String(), user.Email, aDuration, false)

	rDuration, err := time.ParseDuration(as.cfg.RefreshTokenTTL)
	if err != nil {
		as.log.Error("Invalid duration format", slog.String("op", op))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := jwt.NewToken(user.Id.String(), user.Email, rDuration, true)

	if err != nil {
		as.log.Error("Error while creating token: ", slog.String(op, err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrCreatingToken)
	}

	return &authdto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (as *DefaultAuthService) GetAuthedUser(ctx context.Context, t *authdto.Token) (*authdto.AuthedUserResponse, error) {
	const op = "service.GetAuthedUser"

	token, err := t.Validate()
	if err != nil {
		as.log.Error("%s: %w", op, err)
		if errors.Is(err, jwtvalidate.ErrTokenExpired) {
		return nil, ErrTokenExpired
	}
		return nil, ErrValidatingToken
	}

	email, err := token.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrValidatingToken)
	}

	user, err := as.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrNoUserFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s: %w", op, err) 
	}

	return &authdto.AuthedUserResponse{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		RoleName: *user.Role,
		RegistrationDate: user.RegistrationDate,
	}, nil
}
	
func (as *DefaultAuthService) RefreshToken(ctx context.Context, t *authdto.Token) (*authdto.Token, error) {
	const op = "service.RefreshToken"

	//validating a token
	token, err := t.Validate()
	if err != nil {
		as.log.Error("Token validation error", slog.String(op, err.Error()))
		if errors.Is(err, jwtvalidate.ErrTokenExpired) {
		return nil, ErrTokenExpired
	}
		return nil, ErrValidatingToken
	}

	claims, ok := token.Claims.(*jwt.CustomTokenClaims)

	if !ok {
		as.log.Error("Error while retrieving token claims")
		return nil, fmt.Errorf("%s: %w", op, ErrValidatingToken)
	}

	if !claims.IsRefresh {
		as.log.Debug("Token is not refresh.")
		return nil, fmt.Errorf("%s: %w", op, ErrValidatingToken)
	}

	userId := claims.RegisteredClaims.Subject
	dbUser, err := as.repo.FindUserById(ctx, userId)

	if err != nil {
		as.log.Debug("User query error", slog.String(op, err.Error()))
		if errors.Is(err, repo.ErrNoUserFound) {
			as.log.Debug("User was not found", slog.String(op, err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserWasNotFound)
		} else if errors.Is(err, repo.ErrInvalidObjectId) {
			as.log.Debug("Invalid ObjectId format", slog.String(op, err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		as.log.Debug("Error while querying for User", slog.String(op, err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err) 
	}

	duration, err := time.ParseDuration(as.cfg.AccessTokenTTL)
	if err != nil {
		as.log.Error("Invalid duration format", slog.String("op", op))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	accessToken, err := jwt.NewToken(userId, dbUser.Email, duration, false)

	if err != nil {
		as.log.Error("Error while creating token: ", slog.String(op, err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrCreatingToken)
	}
	
	return &authdto.Token{
		Token: accessToken,
	}, nil
}