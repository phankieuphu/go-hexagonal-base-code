package database

import (
	internal_config "account-service/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DatabaseProvider struct {
	db *dynamodb.Client
}

func NewDatabaseProvider(cfg *internal_config.Config, ctx context.Context) (DatabaseProvider, error) {
	dbCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		return DatabaseProvider{}, err
	}
	return DatabaseProvider{
		db: dynamodb.NewFromConfig(dbCfg),
	}, nil
}
