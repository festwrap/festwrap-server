package event

import (
	"context"
	"festwrap/internal/messaging"
	"festwrap/internal/serialization"
	"time"
)

const (
	defaultTimeout = time.Second * 10
)

type PublishEventObserver[T Event] struct {
	publisher  messaging.Publisher
	topic      string
	timeout    time.Duration
	serializer serialization.Serializer[EventWrapper[T]]
}

func NewPublishEventObserver[T Event](publisher messaging.Publisher, topic string) PublishEventObserver[T] {
	serializer := serialization.NewJsonSerializer[EventWrapper[T]]()
	return PublishEventObserver[T]{
		publisher:  publisher,
		topic:      topic,
		serializer: &serializer,
		timeout:    defaultTimeout,
	}
}

func (p *PublishEventObserver[T]) WithTimeout(timeout time.Duration) *PublishEventObserver[T] {
	p.timeout = timeout
	return p
}

func (o PublishEventObserver[T]) Update(playlistEvent EventWrapper[T]) error {
	eventBytes, err := o.serializer.Serialize(playlistEvent)
	if err != nil {
		return err
	}

	// Run in the background so we do not wait for publish confirmation
	go func() {
		timeoutCtx, ctxCancel := context.WithDeadline(context.Background(), time.Now().Add(o.timeout))
		defer ctxCancel()
		o.publisher.Publish(timeoutCtx, o.topic, eventBytes)
	}()
	return nil
}
