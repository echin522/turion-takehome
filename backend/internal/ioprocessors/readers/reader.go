package readers

import "context"

// Reader defines the method needed for an io-reader, aka what comes in
// Readers should not perform any transformation of data. What comes in is just
// bytes and what is passed to the writer is just bytes
type Reader interface {
	Read(context.Context, []byte) (int, error)
}
