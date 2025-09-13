package bus

import (
	// "context"
	"errors"
	// "strings"
	"sync"
	"time"

	// "github.com/google/uuid"
)

var (
	ErrClosed = errors.New("topic closed")
)

type Topic struct {
	name 		string
	buffer		int
	queue		chan Message
	dlq			chan Message
	backoff 	Backoff
	
	mu			sync.RWMutex
	subscribers	map[string]*Sub

	closed 		bool
	wg			sync.WaitGroup
}

type TopicConfig struct {
	Name 	string
	Buffer 	int
	DLQSize	int
	Backoff	Backoff
}

func NewTopic(cfg TopicConfig) *Topic {
	if cfg.Buffer <=0 {
		cfg.Buffer = 1024
	}
	if cfg.DLQSize <0 {
		cfg.DLQSize = 256
	}

	if cfg.Backoff == nil {
		cfg.Backoff = ExpBackoff{Base: 50 * time.Millisecond, Factor: 2.0, Max: 2 * time.Second}
	}
	t := &Topic{
		name: cfg.Name,
		buffer: cfg.Buffer,
		queue: make(chan Message, cfg.Buffer),
		dlq: make(chan Message, cfg.DLQSize),
		backoff: cfg.Backoff,
		subscribers: make(map[string]*Sub),
	}
	return t
}