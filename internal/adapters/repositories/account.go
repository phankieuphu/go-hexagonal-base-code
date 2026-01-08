package repositories

import (
	"account-service/config"
	"account-service/internal/domain/entities"
	"account-service/internal/domain/ports"
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type AccountRepository struct {
	dbService ports.IDatabaseService
	tableName string
}

func NewAccountRepository(dbService ports.IDatabaseService, internalConfig config.Config) (AccountRepository, error) {
	return AccountRepository{dbService: dbService, tableName: internalConfig.Database.AccountTableName}, nil
}

func (e *AccountRepository) Save(ctx context.Context, entry entities.Account) error {
	item, err := attributevalue.MarshalMap(entry)
	if err != nil {
	}
	e.dbService.Insert(ctx, item)
	if err != nil {
		return err
	}
	return nil
}
