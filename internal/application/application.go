package application

import (
	"account-service/config"
	"account-service/internal/adapters/consumer"
	"account-service/internal/adapters/database"
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

	dbProvider, err := database.NewDatabaseProvider(config, ctx)
	if err != nil {
		log.Fatalf("failed to init DATABASE provider: %v", err)
	}
	dbService := database.NewDatabaseService(ctx, dbProvider, config.Database.AccountTableName)

	entryAccountRepository, err := repositories.NewAccountRepository(dbService, *config)
	if err != nil {
		log.Fatalf("failed to init REPOSIORTY provider: %v", err)
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
	accountingConsumer, err := consumer.NewAccountConsumer(ctx, queueProvider, config, entryAccountService, queuURL)
	if err != nil {
		log.Fatalf("failed to init QUEUE consumer: %v", err)
	}
	log.Println("Accont Application Started")
	accountingConsumer.Start(ctx)
}
