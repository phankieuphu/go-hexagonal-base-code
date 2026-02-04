package application

import (
	"account-service/config"
	"account-service/internal/adapters/consumer"
	database_provider "account-service/internal/adapters/database/provider"
	"account-service/internal/adapters/repositories"
	"account-service/internal/domain/services"
	"context"
	"log"

	"github.com/joho/godotenv"
)

type Application struct {
}

func AccountApplication(ctx context.Context) {
	godotenv.Load()
	config := config.LoadConfig()

	database, err := database_provider.NewMySQLClient(*config)
	if err != nil {
		log.Fatalf("Failed to init database service %s", err.Error())
	}

	entryAccountRepository := repositories.NewAccountRepository(database)
	if err != nil {
		log.Fatalf("failed to init REPOSITORY provider: %v", err)
	}

	entryAccountService := services.NewAccountService(*config, entryAccountRepository)

	queueClient, err := consumer.NewSQSClient(*config, ctx)
	if err != nil {
		log.Fatalf("failed to init QUEUE client: %v", err)
	}
	queueProvider, err := consumer.NewQueueProvider(*queueClient)
	if err != nil {
		log.Fatalf("failed to init QUEUE provider: %v", err)
	}
	queuURL := config.SqsTopic.Account
	accountConsumer, err := consumer.NewAccountConsumer(ctx, queueProvider, config, entryAccountService, queuURL)
	if err != nil {
		log.Fatalf("failed to init QUEUE consumer: %v", err)
	}
	log.Println("Accont Application Started")
	accountConsumer.Start(ctx)
}
