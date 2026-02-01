package services

import (
	"account-service/config"
	"account-service/internal/domain/entity"
	"account-service/internal/domain/ports"
	"context"
)

type AccountService struct {
	config     config.Config
	repository ports.AccountRepository
}

// Save implements ports.AccountService.
func (e *AccountService) Save(ctx context.Context, entity entity.Account) error {
	// Handle business here:

	e.repository.Create(ctx, entity)

	return nil
}

func NewAccountService(cfg config.Config, repository ports.AccountRepository) ports.AccountService {
	return &AccountService{
		repository: repository,
		config:     cfg,
	}
}
