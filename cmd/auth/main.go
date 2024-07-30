package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/FelishaK/grpcAuth/internal/app"
	"github.com/FelishaK/grpcAuth/internal/config"
)

func main() {
	// init config
	cfg, err := config.LoadConfig()

	//init logger
	log := setUpLogger()
	
	if err != nil {
		log.Error("Config was not loaded", slog.String("err", err.Error()))
		panic(err)
	}

	log.Info("Starting application . .  .")
	// run server
	srv := app.NewApp(log, cfg)

	go srv.GRPCSrv.RunServer()
	

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop
	log.Info("Received signal", slog.String("signal", signal.String()))

	srv.GRPCSrv.StopServer()
	log.Info("Stopping server . .  .")
}

func setUpLogger() *slog.Logger {
	var log *slog.Logger
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	log = slog.New(slog.NewJSONHandler(os.Stdout, &opts))

	return log
}