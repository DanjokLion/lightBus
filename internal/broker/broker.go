package broker

import (
	"context"
	"sync"
	"time"
)

type Message struct {
	Topic     string
	Payload   []byte
	Retry     int
	Timestamp time.Time
}

type Subscriber chan Message

type Broker interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(topic string) Subscriber
	Ack(msg Message)
	MoveToDLQ(msg Message)
	GetDLQ(topic string) []Message
}

type InMemoryBroker struct {
	mu       sync.RWMutex
	topics   map[string][]Subscriber
	dlq      map[string][]Message
	retryMax int
}

func NewInMemoryBroker() *InMemoryBroker {
	return &InMemoryBroker{
		topics:   make(map[string][]Subscriber),
		dlq:      make(map[string][]Message),
		retryMax: 3,
	}
}

func (b *InMemoryBroker) Publish(ctx context.Context, topic string, data []byte) error {
	b.mu.RLock()
	subs, ok := b.topics[topic]
	b.mu.RUnlock()

	if !ok || len(subs) == 0 {
		return nil
	}

	msg := Message{
		Topic:     topic,
		Payload:   data,
		Retry:     0,
		Timestamp: time.Now(),
	}

	for _, sub := range subs {
		select {
		case sub <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (b *InMemoryBroker) Subscribe(topic string) Subscriber {
	ch := make(Subscriber, 10)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.topics[topic] = append(b.topics[topic], ch)
	return ch
}

func (b *InMemoryBroker) Ack(msg Message) {
}

func (b *InMemoryBroker) MoveToDLQ(msg Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dlq[msg.Topic] = append(b.dlq[msg.Topic], msg)
}

func (b *InMemoryBroker) GetDLQ(topic string) []Message {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dlq[topic]
}
