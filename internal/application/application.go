package application

import (
	"account-service/config"
	"account-service/internal/adapters/cache"
	"account-service/internal/adapters/consumer"
	database_provider "account-service/internal/adapters/database/provider"
	"account-service/internal/adapters/kafka"
	"account-service/internal/adapters/repository"
	"account-service/internal/domain/services"
	"context"
	"log"

	ginhttp "account-service/internal/adapters/http"

	"github.com/joho/godotenv"
)

type Application struct {
}

func AccountApplication(ctx context.Context) {
	godotenv.Load()
	cfg := config.LoadConfig()

	// database
	database, err := database_provider.NewMySQLClient(*cfg)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// repository & service
	entryAccountRepository := repository.NewAccountRepository(database)
	entryAccountService := services.NewAccountService(*cfg, entryAccountRepository)

	// SQS consumer
	queueClient, err := consumer.NewSQSClient(*cfg, ctx)
	if err != nil {
		log.Fatalf("failed to init SQS client: %v", err)
	}
	queueProvider, err := consumer.NewQueueProvider(*queueClient)
	if err != nil {
		log.Fatalf("failed to init queue provider: %v", err)
	}
	accountConsumer, err := consumer.NewAccountConsumer(ctx, queueProvider, cfg, entryAccountService, cfg.SqsTopic.Account)
	if err != nil {
		log.Fatalf("failed to init account consumer: %v", err)
	}

	// Redis cache
	redisCache, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to init Redis: %v", err)
	}
	_ = redisCache

	// Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("failed to init Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka, func(ctx context.Context, key, value []byte) error {
		log.Printf("kafka message received key=%s value=%s", key, value)
		return nil
	})
	if err != nil {
		log.Fatalf("failed to init Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	// HTTP server (gin)
	httpServer := ginhttp.NewServer(cfg.API, entryAccountService)

	log.Println("Account Application Started")

	go accountConsumer.Start(ctx)
	go kafkaConsumer.Start(ctx)
	httpServer.Start()
}
