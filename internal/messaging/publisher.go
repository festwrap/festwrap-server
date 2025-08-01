package messaging

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
}
