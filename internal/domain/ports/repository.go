package ports

import (
	"account-service/internal/domain/entities"
	"context"
)

type AccountRepository interface {
	Create(ctx context.Context, account entities.Account)
}
