// infra/publisher/sqs_publisher.go
package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSPublisher struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSPublisher(queueURL string) (*SQSPublisher, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	if queueURL == "" {
		return nil, fmt.Errorf("queueURL não pode ser vazio")
	}

	return &SQSPublisher{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (p *SQSPublisher) Publish(eventType string, payload interface{}) error {
	msg := map[string]interface{}{
		"event_type": eventType,
		"data":       payload,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: aws.String(string(body)),
	})
	return err
}
