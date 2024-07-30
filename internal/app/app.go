package app

import (
	"fmt"
	"log/slog"

	authapp "github.com/FelishaK/grpcAuth/internal/app/auth"
	"github.com/FelishaK/grpcAuth/internal/config"
	repo "github.com/FelishaK/grpcAuth/internal/repository"
	"github.com/FelishaK/grpcAuth/internal/service"
)

type App struct {
	GRPCSrv *authapp.Server
	Config *config.Config
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	const op = "app.NewApp"
	storage, err := repo.NewStorage(config.Database.MongoUser, config.Database.MongoPassword, config.Database.MongoHostname, config.Database.MongoDBName, config.Database.MongoPort, config.Database.Timeout)

	if err != nil {
		log.Error("Error while initializing db", slog.String(op, err.Error()))
		panic(err)
	}
	fmt.Printf("CONFIG: %+v", config)
	authService := service.NewAuthService(log, config, storage)
	srv := authapp.NewApp(log, config, authService)
	return &App{
		GRPCSrv: srv,
	}
}