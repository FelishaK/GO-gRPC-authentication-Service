package authapp

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/FelishaK/grpcAuth/internal/config"
	authroutes "github.com/FelishaK/grpcAuth/internal/grpcRoutes/authRoutes/auth"

	// authroutes "github.com/FelishaK/grpcAuth/internal/grpcRoutes/authroutes"
	"github.com/FelishaK/grpcAuth/internal/service"
	"google.golang.org/grpc"
)

type Server struct {
	log *slog.Logger
	gRPCServer *grpc.Server
	hostname string
	port int
}

func NewApp(log *slog.Logger, config *config.Config, authService service.AuthService) *Server {
	grpcServer := grpc.NewServer()
	authroutes.Register(grpcServer, authService)

	return &Server{
		log: log,
		gRPCServer: grpcServer,
		hostname: config.Grpc.GrpcHost,
		port: config.Grpc.Port,	
	}
}


func (a *Server) RunServer() error {
	const op = "auth.RunServer"

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.hostname, a.port))

	if err != nil {
		a.log.Error("failed to listen", slog.String("op: ", op), slog.String("localhost",  a.hostname), slog.String("port: ", strconv.Itoa(a.port)))
		panic(err)
	}

	if err :=  a.gRPCServer.Serve(lis); err != nil {
		a.log.Error("failed to serve", slog.String("op: ", op))
		panic(err)
	}
	a.log.Info("Server has been launched successfully")
	return nil
}

// Gracefully stop the server
func (a *Server) StopServer() {
	const op = "auth.StopServer" 

	a.gRPCServer.GracefulStop()
	a.log.Info("Stopping the server", slog.String("op", op))
}