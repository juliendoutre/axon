package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juliendoutre/axon/internal/config"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var (
	Semver        string
	GitCommitHash string
	BuildTime     string
	GoVersion     string
	Os            string
	Arch          string
)

func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	defer func() { _ = logger.Sync() }()

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	pgPool, err := pgxpool.New(ctx, config.PostgresURL().String())
	if err != nil {
		logger.Panic("Connecting to DB", zap.Error(err))
	}
	defer pgPool.Close()

	temporalClient, err := client.Dial(client.Options{
		HostPort: net.JoinHostPort(os.Getenv("TEMPORAL_HOST"), os.Getenv("TEMPORAL_PORT")),
	})
	if err != nil {
		logger.Panic("Creating Temporal client", zap.Error(err))
	}
	defer temporalClient.Close()

	temporalWorker := worker.New(temporalClient, os.Getenv("TEMPORAL_TASK_QUEUE"), worker.Options{})

	if err := temporalWorker.Run(worker.InterruptCh()); err != nil {
		logger.Panic("Starting Temporal worker", zap.Error(err))
	}
}
