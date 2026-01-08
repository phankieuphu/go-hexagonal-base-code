package services

import (
	"account-service/config"
	"account-service/internal/adapters/repositories"
	"account-service/internal/domain/entities"
	"account-service/internal/domain/ports"
	"context"

	"github.com/google/uuid"
)

type AccountService struct {
	config     config.Config
	repository repositories.AccountRepository
}

// Save implements ports.AccountService.
func (e *AccountService) Save(ctx context.Context, input string) {
	// 2. Map to command
	data := entities.Account{
		ID: uuid.NewString(),
	}
	err := e.repository.Save(ctx, data)
	if err != nil {

	}

}

func NewAccountService(cfg config.Config, repository repositories.AccountRepository) ports.AccountService {
	return &AccountService{
		repository: repository,
		config:     cfg,
	}
}
