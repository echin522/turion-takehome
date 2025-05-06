package readers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ChannelReader struct {
	logger  *zap.Logger
	channel <-chan []byte
}

func NewChannelReader(logger *zap.Logger, byteChannel <-chan []byte) *ChannelReader {
	return &ChannelReader{
		logger:  logger,
		channel: byteChannel,
	}
}

// Read implements the ContextReader interface used on processors.Processors.
// Use this over other processors to keep state of your internal process.
func (cr *ChannelReader) Read(ctx context.Context, b []byte) (n int, err error) {
	readChannelCtx, cancel := context.WithTimeout(ctx, time.Second*5)

	defer cancel()
	select {
	case <-readChannelCtx.Done():
		return 0, nil
	case data, ok := <-cr.channel:
		if !ok {
			return 0, nil
		}

		if len(b) > len(data) {
			b = b[:len(data)]
		}

		copy(b, data)
		cr.logger.Debug("Reading message from channel")

		return len(data), nil
	}
}
