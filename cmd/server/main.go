package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

//nolint:gochecknoglobals
var (
	Semver        string
	GitCommitHash string
	BuildTime     string
	GoVersion     string
	Os            string //nolint:varnamelen
	Arch          string
)

//nolint:funlen,cyclop
func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	defer func() { _ = logger.Sync() }()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	creds, err := credentials.NewServerTLSFromFile(os.Getenv("SERVER_CERT_PATH"), os.Getenv("SERVER_KEY_PATH"))
	if err != nil {
		logger.Panic("Loading TLS credentials", zap.Error(err))
	}

	jwkStore, err := keyfunc.NewDefaultCtx(ctx, strings.Split(os.Getenv("SERVER_JWKS_URLS"), ","))
	if err != nil {
		logger.Panic("Starting JWKs store", zap.Error(err))
	}

	policy, err := os.ReadFile(os.Getenv("SERVER_POLICY_PATH"))
	if err != nil {
		logger.Panic("Loading policy", zap.Error(err))
	}

	authenticator := &server.Authenticator{
		Parser: jwt.NewParser(),
		Store:  jwkStore.KeyfuncCtx(ctx),
		Policy: string(policy),
	}

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			authenticator.Authenticate,
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	grpcServer := grpc.NewServer(grpcOptions...)

	pgPool, err := pgxpool.New(ctx, config.PostgresURL().String())
	if err != nil {
		logger.Panic("Connecting to DB", zap.Error(err))
	}
	defer pgPool.Close()

	parsedBuildTime, err := time.Parse(time.RFC3339, BuildTime)
	if err != nil {
		logger.Panic("Parsing build time", zap.Error(err))
	}

	server, err := server.New(
		&v1.Version{
			Semver:        Semver,
			GitCommitHash: GitCommitHash,
			BuildTime:     timestamppb.New(parsedBuildTime),
			GoVersion:     GoVersion,
			Os:            Os,
			Arch:          Arch,
		},
		pgPool,
	)
	if err != nil {
		logger.Panic("Creating server", zap.Error(err))
	}

	healthgrpc.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)
	v1.RegisterAxonServer(grpcServer, server)

	grpcPort, err := strconv.ParseInt(os.Getenv("SERVER_PORT"), 10, 64)
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
