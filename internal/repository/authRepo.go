package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	mongoapp "github.com/FelishaK/grpcAuth/internal/app/mongo"
	"github.com/FelishaK/grpcAuth/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserExists = errors.New("user already exists")
	ErrNoUserFound = errors.New("user was not found")
	ErrInvalidObjectId = errors.New("invalid ObjectId format")
)

type AuthRepository interface {
	AddUser(ctx context.Context, user domain.User) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserById(ctx context.Context, id string) (*domain.User, error)
}

type Storage struct {
	mongo *mongoapp.MongoDB
}

func NewStorage(mongoUser string, mongoPassword string, mongoHost string, mongoDB string, mongoPort int, timeout time.Duration) (*Storage, error) {
	const op = "repo.NewStorage"

	store, err := mongoapp.NewMongoDB(mongoUser, mongoPassword, mongoHost, mongoDB, mongoPort, timeout)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		mongo: store,
	}, nil
}

func (s *Storage) AddUser(ctx context.Context, user domain.User) (string, error) {
	const op = "repo.AddUser"
	res, err := s.mongo.UsersCol.InsertOne(ctx, user)
	
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	userId := res.InsertedID
	return userId.(primitive.ObjectID).Hex(), nil
}

func (s *Storage) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "repo.FindUserByEmail"

	filter := bson.M{"email": email}
	var user domain.User

	err := s.mongo.UsersCol.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &domain.User{}, fmt.Errorf("%s: %w", op, ErrNoUserFound)
		}
		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (s *Storage) FindUserById(ctx context.Context, id string) (*domain.User, error) {
	const op = "repo.FindUserById"

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, "Invalid ObjectId")
	}
	fmt.Println("\nID: ", objId)
	filter := bson.M{"_id": objId}
	var user domain.User

	err = s.mongo.UsersCol.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &domain.User{}, fmt.Errorf("%s: %w", op, ErrNoUserFound)
		}
		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}