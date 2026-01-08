package ports

import (
	"context"
)

type AccountService interface {
	Save(context.Context, string)
}
