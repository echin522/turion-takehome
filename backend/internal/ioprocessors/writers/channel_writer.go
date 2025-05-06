package writers

import (
	"context"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
)

type ChannelWriter struct {
	logger *zap.Logger

	close   atomic.Bool
	mu      *sync.Mutex
	channel chan<- []byte
}

func NewChannelWriter(logger *zap.Logger, channel chan<- []byte) *ChannelWriter {
	return &ChannelWriter{
		logger:  logger,
		channel: channel,
		mu:      &sync.Mutex{},
	}
}

// Writer the bytes to a channel.
func (w *ChannelWriter) Write(ctx context.Context, b []byte) (int, error) {
	w.channel <- b

	return len(b), nil
}

func (w *ChannelWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.close.Load() {
		return nil
	}

	w.close.Store(true)

	close(w.channel)
	return nil
}
