package database

import (
	"account-service/internal/domain/ports"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// type Database struct {
// 	provider  DatabaseProvider
// 	tableName string
// }

type DatabaseService struct {
	provider  DatabaseProvider
	ctx       context.Context
	tableName string
}

// Close implements ports.IDatabaseService.
func (d DatabaseService) Close() {

}

// Delete implements ports.IDatabaseService.
func (d DatabaseService) Delete() {
	panic("unimplemented")
}

// GetList implements ports.IDatabaseService.
func (d DatabaseService) GetList() {
	fmt.Println("GET LIST")
}

// Insert implements ports.IDatabaseService.
func (d DatabaseService) Insert(ctx context.Context, items any) {
	av, err := attributevalue.MarshalMap(items)
	if err != nil {
		return
	}
	input := &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      av,
	}
	d.provider.db.PutItem(ctx, input)
	fmt.Println("INSERT", items)

}

// Ping implements ports.IDatabaseService.
func (d DatabaseService) Ping() {
	panic("unimplemented")
}

// Update implements ports.IDatabaseService.
func (d DatabaseService) Update() {
	panic("unimplemented")
}

func NewDatabaseService(ctx context.Context, provider DatabaseProvider, tableName string) ports.IDatabaseService {
	return DatabaseService{
		provider:  provider,
		ctx:       ctx,
		tableName: tableName,
	}
}
