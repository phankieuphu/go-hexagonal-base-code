package consumer

import (
	internal_config "account-service/config"
	"account-service/internal/domain/ports"
	"context"
	"encoding/json"
	"log"
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

func (a *AccountConsumer) Start(ctx context.Context) {

	queueURL := a.queueURL

	log.Println("SQS consumer started...")

	for {
		resp, err := a.provider.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 5,
			WaitTimeSeconds:     20, // long polling
			VisibilityTimeout:   30,
		})
		if err != nil {
			log.Println("receive message error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(resp.Messages) == 0 {
			continue
		}

		for _, msg := range resp.Messages {
			err := a.ProcessMessage(ctx, msg)
			if err != nil {
				log.Println("processing failed:", err)
				continue
			}
		}
	}
}

func (a *AccountConsumer) ProcessMessage(
	ctx context.Context,
	msg types.Message,
) error {
	log.Println("📩 raw SQS message received")

	// 1. Unwrap SNS envelope
	var snsMsg ports.SNSMessage
	if err := json.Unmarshal([]byte(*msg.Body), &snsMsg); err != nil {
		return err
	}

	a.entryAccountService.Save(ctx, snsMsg.Message)

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
		log.Println("message deleted")
	}

	return err
}
