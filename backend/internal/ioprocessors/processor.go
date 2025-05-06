package ioprocessors

import (
	"context"
	"errors"
	"sync"
	"turion-takehome/internal/ioprocessors/quarantiners"
	"turion-takehome/internal/ioprocessors/readers"
	"turion-takehome/internal/ioprocessors/writers"

	"go.uber.org/zap"
)

// An arbitrary soft limit for buffer size. Telemetry messages are pretty small,
// especially coming off the radio, so it would be very surprising if this limit
// needs to be broken
const BUFFER_SIZE_WARNING_LIMIT = 2048

type Processor struct {
	logger *zap.Logger

	pool        *sync.Pool
	reader      readers.Reader
	writer      writers.Writer
	quarantiner quarantiners.Quarantiner

	option *ProcessorOptions
}

// NewProcessor will run a process with two different a reader and a writer.
// Context will be passed along to both.
func NewProcessor(
	logger *zap.Logger,
	bufferSize int,
	reader readers.Reader,
	writer writers.Writer,
	quarantiner quarantiners.Quarantiner,
	opts ...ProcessorOption,
) *Processor {
	if bufferSize > BUFFER_SIZE_WARNING_LIMIT {
		logger.Warn(
			"Buffer size exceeds soft limit. Reconsider if this large buffer is necessary",
			zap.Int("Buffer size", bufferSize),
			zap.Int("Soft limit", BUFFER_SIZE_WARNING_LIMIT))
	}

	options := &ProcessorOptions{}
	for _, opt := range opts {
		opt(options)
	}

	pool := sync.Pool{
		New: func() any {
			return make([]byte, bufferSize)
		},
	}

	return &Processor{
		logger:      logger,
		pool:        &pool,
		quarantiner: quarantiner,
		reader:      reader,
		writer:      writer,
		option:      options,
	}
}

// Start will initialize the processor
func (p *Processor) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			ctxErr := ctx.Err()
			if p.option.returnContextError {
				if ctxErr != nil {
					p.logger.Error("Context has completed, processor stopped", zap.Error(ctxErr))
				}

				return ctxErr
			}

			return nil

		default:
			// listen to reader
			interfaceFromPool := p.pool.Get()
			if interfaceFromPool == nil {
				return nil
			}

			bytesSlice, ok := interfaceFromPool.([]byte)
			if !ok {
				return errors.New("could not get bytes from pool")
			}

			bytesRead, err := p.reader.Read(ctx, bytesSlice)
			if err != nil {
				return err
			}

			if bytesRead == 0 {
				continue
			}

			bytesToSend := bytesSlice[:bytesRead]
			p.pool.Put(bytesSlice)

			_, err = p.writer.Write(ctx, bytesToSend)
			if err != nil {

				handleWriteError := p.handleWriteErrorToQuarantine(ctx, bytesToSend, err)
				if handleWriteError != nil {
					err = handleWriteError
				}

				// return as an option
				if p.option.returnOnWriterError {
					return err
				}
			}

			// TODO: Do metrics here if you have time
		}
	}
}

func (p *Processor) Close() error {
	err := p.writer.Close()
	if err != nil {
		p.logger.Error("error closing writer on processor", zap.Error(err))
	}

	return err
}

func (p *Processor) handleWriteErrorToQuarantine(ctx context.Context, bytesToSend []byte, writeError error) error {
	if p.option.ignoreContextError && (errors.Is(writeError, context.Canceled) || errors.Is(writeError, context.DeadlineExceeded)) {
		return nil
	}

	_, quarantineError := p.quarantiner.Quarantine(ctx, bytesToSend, writeError)
	if quarantineError != nil {
		err := errors.Join(writeError, quarantineError)
		return err
	}

	return nil
}

type ProcessorOptions struct {
	returnOnWriterError bool
	returnContextError  bool
	ignoreContextError  bool
}

type ProcessorOption func(*ProcessorOptions)

// WithReturnOnWriterError will exit the process when
// we cannot get the message.
// It will return after we send the message to quarantine.
func WithReturnOnWriterError() ProcessorOption {
	return func(o *ProcessorOptions) {
		o.returnOnWriterError = true
	}
}

// WithReturnContextError will exit the process when
// and return context when it logs an error.
// It will return an error.
func WithReturnContextError() ProcessorOption {
	return func(o *ProcessorOptions) {
		o.returnOnWriterError = true
	}
}

// WithQuarantineIgnoreContextError will skip the quarantine
// method if set and error from Writer is Context.Cancelled or Context.Deadline
// This should be used sparingly, and intended for testing parallel tests due to
// timeouts when run concurrently.
func WithQuarantineIgnoreContextError() ProcessorOption {
	return func(o *ProcessorOptions) {
		o.ignoreContextError = true
	}
}
