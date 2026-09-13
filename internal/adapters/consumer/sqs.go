package consumer

import (
	internal_config "account-service/config"
	"account-service/pkg/logger"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewSQSClient(internal_config internal_config.Config, ctx context.Context) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(internal_config.Region))
	if err != nil {
		logger.Error("load aws config failed", "error", err)
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}
