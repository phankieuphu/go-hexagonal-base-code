package consumer

import (
	internal_config "account-service/config"
	"account-service/internal/domain/ports"
	"account-service/pkg/logger"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type AccountConsumer struct {
	cfg                 *internal_config.Config
	entryAccountService ports.AccountService
	queueURL            string
	provider            QueueProvider
}

func NewAccountConsumer(ctx context.Context, provider QueueProvider, cfg *internal_config.Config, entryAccountService ports.AccountService, queueURL string) (ports.IConsumer, error) {
	return &AccountConsumer{
		cfg:                 cfg,
		entryAccountService: entryAccountService,
		queueURL:            queueURL,
		provider:            provider,
	}, nil

}

// Start long-polls the queue until ctx is cancelled. A message already
// received is processed to completion before Start returns.
func (a *AccountConsumer) Start(ctx context.Context) {
	queueURL := a.queueURL
	if queueURL == "" {
		logger.Warn("SQS queue URL not configured, account consumer disabled")
		return
	}

	logger.Info("SQS consumer started", "queue", queueURL)

	for ctx.Err() == nil {
		resp, err := a.provider.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     20, // long polling
			VisibilityTimeout:   30,
		})
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Error("receive message error", "error", err)
			select {
			case <-ctx.Done():
			case <-time.After(2 * time.Second):
			}
			continue
		}

		for _, msg := range resp.Messages {
			// Detach from ctx so a shutdown mid-batch doesn't abort the
			// delete of a message that was already processed.
			if err := a.ProcessMessage(context.WithoutCancel(ctx), msg); err != nil {
				logger.Error("processing failed", "error", err)
			}
		}
	}

	logger.Info("SQS consumer stopped")
}

func (a *AccountConsumer) ProcessMessage(
	ctx context.Context,
	msg types.Message,
) error {
	logger.Info("raw SQS message received")

	// 1. Unwrap SNS envelope

	//if err := json.Unmarshal([]byte(*msg.Body), &snsMsg); err != nil {
	//	return err
	//}
	//
	//a.entryAccountService.Save(ctx, snsMsg.Message)

	// 5. Delete message after success
	return a.DeleteMessage(ctx, *msg.ReceiptHandle)
}

func (a *AccountConsumer) DeleteMessage(
	ctx context.Context,
	receiptHandle string,
) error {
	queueURL := a.queueURL
	_, err := a.provider.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	})

	if err == nil {
		logger.Info("message deleted")
	}

	return err
}
