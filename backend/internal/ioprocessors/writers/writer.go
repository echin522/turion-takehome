package writers

import (
	"context"
	"io"
)

// Writer defines the methods needed for an io-writer, aka what goes out
// Writers are where you will perform any data transformation e.g. bytes to
// tleemetry message and telemetry message to SQL
type Writer interface {
	Write(context.Context, []byte) (int, error)
	io.Closer
}
