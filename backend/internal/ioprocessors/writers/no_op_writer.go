package writers

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type NoOpWriter struct {
	logger *zap.Logger
}

func NewNoOpWriter(
	logger *zap.Logger,
) (*NoOpWriter, error) {
	var errs error
	if logger == nil {
		errs = errors.Join(errors.New("logger cannot be nil"))
	}

	if errs != nil {
		return nil, errs
	}

	return &NoOpWriter{
		logger: logger,
	}, nil
}

func (w *NoOpWriter) Write(ctx context.Context, b []byte) (int, error) {
	w.logger.Info("No write occurred, but processor is working")

	return len(b), nil
}

func (w *NoOpWriter) Close() error {
	return nil
}
