package bus

import (
	"context"
	"time"
)

type MessageID string

type Message struct {
	ID	MessageID
	Topic string
	Key string
	Payload []byte
	Meta map[string]string
	Attempt int
	MaxRetry int
	CreatedAt time.Time
}

type Handler func (ctx context.Context, msg Message) error

type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	if e.Err == nil {
		return "retryable"
	}
	return e.Err.Error()
}

func Retry(err error) error {
	return &RetryableError{Err: err}
}