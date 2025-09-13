package bus

import "context"

type SubsOptions struct {
	Workers		int
	MaxRetry	int
}

type Sub struct {
	topic	*Topic
	name	string
	handler Handler
	options	SubsOptions
	cancel	context.CancelFunc
	
}