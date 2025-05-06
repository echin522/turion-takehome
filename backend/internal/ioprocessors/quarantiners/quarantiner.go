package quarantiners

import "context"

// Quarantiner is used in the processor encounters an error while reading or
// writing a message. Rather than kill the whole process, we quarantine the
// message elsewhere for review
type Quarantiner interface {
	Quarantine(context.Context, []byte, error) (int, error)
}
