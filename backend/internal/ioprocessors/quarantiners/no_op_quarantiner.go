package quarantiners

import (
	"context"

	"go.uber.org/zap"
)

// NoOpQuarantiner simply logs the message that has been quarantined.
// Quarantined messages include messages with malformed structures. This means
// anomalies are NOT included in messages to be quarantineed
type NoOpQuarantiner struct {
	logger        *zap.Logger
}

func NewNoOpQuarantiner(logger *zap.Logger) (*NoOpQuarantiner, error) {
	return &NoOpQuarantiner{logger: logger}, nil
}

// Quarantine writes a message to a Quarantine Channel for further processing of
// an io message. It returns the number of bytes written to the Quarantine
// Channel and any errors that occurred while writing to it.
func (q *NoOpQuarantiner) Quarantine(ctx context.Context, b []byte, err error) (int, error) {
	q.logger.Warn(
		"Quarantining message but performing no action with it",
		zap.Binary("Bytes", b),
		zap.Error(err))

	return len(b), nil
}
