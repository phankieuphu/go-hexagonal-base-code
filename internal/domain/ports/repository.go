package ports

import (
	"account-service/internal/domain/entity"
	"context"
)

type AccountRepository interface {
	Create(ctx context.Context, account entity.Account)
}
