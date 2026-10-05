package application

import (
	"account-service/config"
	"account-service/internal/adapters/cache"
	"account-service/internal/adapters/consumer"
	database_provider "account-service/internal/adapters/database/provider"
	"account-service/internal/adapters/kafka"
	"account-service/internal/adapters/repository"
	"account-service/internal/domain/services"
	"account-service/pkg/logger"
	"context"
	"sync"
	"time"

	ginhttp "account-service/internal/adapters/http"

	"github.com/joho/godotenv"
	health "github.com/phankieuphu/go-health-check"
)

// shutdownTimeout bounds how long in-flight HTTP requests get to finish
// once a shutdown signal arrives, before the server is forcibly closed.
const shutdownTimeout = 15 * time.Second

// healthCheckTimeout is the default per-check deadline for readiness probes.
const healthCheckTimeout = 2 * time.Second

type Application struct {
}

func AccountApplication(ctx context.Context) {
	godotenv.Load()
	cfg := config.LoadConfig()
	logger.Init(cfg.Logger.ToOptions())

	// database
	database, err := database_provider.NewMySQLClient(*cfg)
	if err != nil {
		logger.Fatal("failed to init database", "error", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		logger.Fatal("failed to get database handle", "error", err)
	}

	// Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logger.Fatal("failed to init Kafka producer", "error", err)
	}

	// Redis cache
	redisCache, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		logger.Fatal("failed to init Redis", "error", err)
	}

	// repository & service
	entryAccountRepository := repository.NewAccountRepository(database)
	entryAccountService := services.NewAccountService(*cfg, entryAccountRepository)

	// SQS consumer
	queueClient, err := consumer.NewSQSClient(*cfg, ctx)
	if err != nil {
		logger.Fatal("failed to init SQS client", "error", err)
	}
	queueProvider, err := consumer.NewQueueProvider(*queueClient)
	if err != nil {
		logger.Fatal("failed to init queue provider", "error", err)
	}
	accountConsumer, err := consumer.NewAccountConsumer(ctx, queueProvider, cfg, entryAccountService, cfg.SqsTopic.Account)
	if err != nil {
		logger.Fatal("failed to init account consumer", "error", err)
	}

	// Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka, func(ctx context.Context, key, value []byte) error {
		logger.Info("kafka message received", "key", string(key), "value", string(value))
		return nil
	})
	if err != nil {
		logger.Fatal("failed to init Kafka consumer", "error", err)
	}

	// HTTP server (gin) with liveness/readiness probes. Only the database is
	// critical: without it the service can't serve any request.
	h := health.New(health.WithTimeout(healthCheckTimeout), health.WithChecks(
		health.Check{Name: "mysql", Critical: true, Func: sqlDB.PingContext},
		health.Check{Name: "kafka", Func: kafkaProducer.Ping},
		health.Check{Name: "redis", Func: redisCache.Ping},
	))
	httpServer := ginhttp.NewServer(cfg.API, entryAccountService, h)

	// Background workers run on their own context, cancelled only after the
	// HTTP server has drained, so in-flight requests can still rely on them.
	workerCtx, stopWorkers := context.WithCancel(context.WithoutCancel(ctx))
	defer stopWorkers()
	var workers sync.WaitGroup

	workers.Go(func() { accountConsumer.Start(workerCtx) })
	workers.Go(func() { kafkaConsumer.Start(workerCtx) })

	serverErr := make(chan error, 1)
	go func() { serverErr <- httpServer.Start() }()

	logger.Info("Account Application Started")

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received, shutting down")
	case err := <-serverErr:
		if err != nil {
			logger.Error("HTTP server error, shutting down", "error", err)
		}
	}

	// 1. Stop accepting requests and let in-flight ones finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown", "error", err)
	}

	// 2. Stop the consumers and wait for any in-progress message to finish
	//    before closing what they use.
	stopWorkers()
	workers.Wait()

	// 3. Release connections.
	if err := kafkaConsumer.Close(); err != nil {
		logger.Error("close Kafka consumer", "error", err)
	}
	if err := kafkaProducer.Close(); err != nil {
		logger.Error("close Kafka producer", "error", err)
	}
	if err := redisCache.Close(); err != nil {
		logger.Error("close Redis", "error", err)
	}
	if err := sqlDB.Close(); err != nil {
		logger.Error("close database", "error", err)
	}

	logger.Info("Account Application stopped")
}
