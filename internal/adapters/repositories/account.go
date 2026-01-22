package repositories

import (
	"account-service/internal/adapters/database/models"
	"account-service/internal/domain/entities"
	"account-service/internal/domain/ports"
	"context"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

// Create implements ports.AccountRepository.
func (a AccountRepository) Create(ctx context.Context, account entities.Account) {
	models := a.toModels(account)
	a.db.Save(models)
	panic("unimplemented")
}

func (a AccountRepository) toModels(account entities.Account) models.Account {
	panic("unimplemented")
}

func (a AccountRepository) toDomain(model models.Account) entities.Account {
	panic("unimplemented")
}

func NewAccountRepository(db *gorm.DB) ports.AccountRepository {
	return AccountRepository{
		db: db,
	}
}
