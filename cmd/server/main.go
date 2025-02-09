package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juliendoutre/axon/internal/config"
	"github.com/juliendoutre/axon/internal/server"
	v1 "github.com/juliendoutre/axon/pkg/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

//nolint:gochecknoglobals
var (
	GoVersion string
	Os        string //nolint:varnamelen
	Arch      string
)

//nolint:funlen
func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	defer func() { _ = logger.Sync() }()

	creds, err := credentials.NewServerTLSFromFile("/etc/axon/server.crt.pem", "/etc/axon/server.key.pem")
	if err != nil {
		logger.Panic("Loading TLS credentials", zap.Error(err))
	}

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	}

	grpcServer := grpc.NewServer(grpcOptions...)

	healthgrpc.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)

	pgURL, err := config.PostgresURL()
	if err != nil {
		logger.Panic("Reading PostgresQL config", zap.Error(err))
	}

	ctx := context.Background()

	pgPool, err := pgxpool.New(ctx, pgURL.String())
	if err != nil {
		logger.Panic("Connecting to DB", zap.Error(err))
	}
	defer pgPool.Close()

	server, err := server.New(&v1.Version{
		GoVersion: GoVersion,
		Os:        Os,
		Arch:      Arch,
	}, pgPool)
	if err != nil {
		logger.Panic("Creating server", zap.Error(err))
	}

	v1.RegisterAxonServer(grpcServer, server)

	grpcPort, err := strconv.ParseInt(os.Getenv("GRPC_PORT"), 10, 64)
	if err != nil {
		logger.Panic("Parsing gRPC port", zap.Error(err))
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Panic("Creating a TCP listener", zap.Error(err))
	}

	go handleSignals(logger, grpcServer)

	logger.Info("Starting the axon server...", zap.Int64("port", grpcPort))

	if err := grpcServer.Serve(listener); err != nil {
		logger.Panic("Serving gRPC request", zap.Error(err))
	}
}

func handleSignals(logger *zap.Logger, server *grpc.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for signal := range signals {
		logger.Warn("Caught a cancellation signal, terminating...", zap.String("signal", signal.String()))
		server.GracefulStop()

		return
	}
}
