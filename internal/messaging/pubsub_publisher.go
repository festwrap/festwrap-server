package messaging

import (
	"context"
	"fmt"

	"festwrap/internal/logging"

	"cloud.google.com/go/pubsub"
)

type PubsubPublisher struct {
	logger logging.Logger
	client *pubsub.Client
}

func NewPubsubPublisher(client *pubsub.Client, logger logging.Logger) PubsubPublisher {
	return PubsubPublisher{client: client, logger: logger}
}

func (p PubsubPublisher) Publish(ctx context.Context, topic string, message []byte) error {
	pubsubTopic := p.client.Topic(topic)
	result := pubsubTopic.Publish(ctx, &pubsub.Message{
		Data: message,
	})

	_, err := result.Get(ctx)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Failed to publish message to topic %s: %v", topic, err))
	}

	return err
}
