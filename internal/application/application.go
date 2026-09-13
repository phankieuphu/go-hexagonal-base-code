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

	ginhttp "account-service/internal/adapters/http"

	"github.com/joho/godotenv"
)

type Application struct {
}

func AccountApplication(ctx context.Context) {
	godotenv.Load()
	cfg := config.LoadConfig()

	logger.Init(logger.Options{
		Level:     cfg.Logger.Level,
		Format:    logger.Format(cfg.Logger.Format),
		AddSource: cfg.Logger.AddSource,
	})

	// database
	database, err := database_provider.NewMySQLClient(*cfg)
	if err != nil {
		logger.Fatal("failed to init database", "error", err)
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

	// Redis cache
	redisCache, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		logger.Fatal("failed to init Redis", "error", err)
	}
	_ = redisCache

	// Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logger.Fatal("failed to init Kafka producer", "error", err)
	}
	defer kafkaProducer.Close()

	// Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka, func(ctx context.Context, key, value []byte) error {
		logger.Info("kafka message received", "key", string(key), "value", string(value))
		return nil
	})
	if err != nil {
		logger.Fatal("failed to init Kafka consumer", "error", err)
	}
	defer kafkaConsumer.Close()

	// HTTP server (gin)
	httpServer := ginhttp.NewServer(cfg.API, entryAccountService)

	logger.Info("Account Application Started")

	go accountConsumer.Start(ctx)
	go kafkaConsumer.Start(ctx)
	httpServer.Start()
}
