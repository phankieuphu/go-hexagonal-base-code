package main

import (
	"account-service/internal/application"
	"context"
)

func main() {
	ctx := context.Background()
	application.AccountApplication(ctx)

}
