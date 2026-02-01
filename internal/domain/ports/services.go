package ports

import (
	"account-service/internal/domain/entity"
	"context"
)

type AccountService interface {
	Save(context.Context, entity.Account) error
}
