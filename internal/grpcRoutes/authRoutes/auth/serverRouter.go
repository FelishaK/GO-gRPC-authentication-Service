package authroutes

import (
	"context"
	"errors"

	auth1 "github.com/FelishaK/authProto/gen/auth"
	authdto "github.com/FelishaK/grpcAuth/internal/dto/auth"
	"github.com/FelishaK/grpcAuth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authService service.AuthService
	auth1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server, authService service.AuthService) {
	auth1.RegisterAuthServer(gRPC, &serverAPI{authService: authService})
}

func (s *serverAPI) Register(ctx context.Context, req *auth1.RegisterRequest) (*auth1.RegisterResponse, error) {
	registerDTO := authdto.ToDTO(req)

	cont := context.Background()
	userId, err := s.authService.Register(cont, registerDTO)

	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		} else if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		}
	}
	return &auth1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *auth1.LoginRequest) (*auth1.LoginResponse, error) {
	loginDTO := authdto.LoginReqToDTO(req)
	token, err := s.authService.Login(ctx, loginDTO)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		} else if errors.Is(err, service.ErrUserWasNotFound) {
			return nil, status.Error(codes.NotFound, "User was not found")
		} else {
			return nil, status.Error(codes.Canceled, "Something went wrong")
		}

	}
	return &auth1.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *serverAPI) GetAuthedUser(ctx context.Context, req *auth1.TokenRequest) (*auth1.AuthedUserResponse, error) {

	tokenDTO := authdto.TokenReqToDTO(req)
	user, err := s.authService.GetAuthedUser(ctx, tokenDTO)

	if err != nil {
		if errors.Is(err, service.ErrUserWasNotFound) {
			return nil, status.Error(codes.NotFound, "User was not found")
		} else if errors.Is(err, service.ErrTokenExpired) {
			return nil, status.Error(codes.PermissionDenied, "Token has been expired")
		} else if errors.Is(err, service.ErrValidatingToken) {
			return nil, status.Error(codes.PermissionDenied, "Invalid Token")
		}
	}

	return &auth1.AuthedUserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		RoleName:  user.RoleName,
	}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *auth1.TokenRequest) (*auth1.RefreshResponse, error) {
	tokenDTO := authdto.TokenReqToDTO(req)
	t, err := s.authService.RefreshToken(ctx, tokenDTO)

	if err != nil {
		if errors.Is(err, service.ErrUserWasNotFound) {
			return nil, status.Error(codes.NotFound, "User was not found")
		} else if errors.Is(err, service.ErrTokenExpired) {
			return nil, status.Error(codes.PermissionDenied, "Token has been expired")
		} else if errors.Is(err, service.ErrValidatingToken) {
			return nil, status.Error(codes.PermissionDenied, "Invalid Token")
		} else {
			return nil, status.Error(codes.Canceled, "Server error")
		}
	}

	return &auth1.RefreshResponse{
		AccessToken: t.Token,
	}, nil
}
