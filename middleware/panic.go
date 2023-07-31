package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// PanicErrorHandler defines a function to call when a panic occurs.
type PanicErrorHandler func(req *request.RequestContext, err error, stack []byte) error

// PanicOptions defines the behavior on how to deal with panics raised by request handlers.
type PanicOptions struct {
	// StackSize establishes the maximum stack buffer to print in bytes.
	StackSize int `json:"stackSize,omitempty"`

	// IncludeAllGoRoutines, if true, then the stack of all the go routines are included.
	IncludeAllGoRoutines bool `json:"includeAllGoRoutines,omitempty"`

	// PanicErrorHandler is an optional custom callback to call if a panic is raised.
	PanicErrorHandler PanicErrorHandler
}

// -----------------------------------------------------------------------------

// NewPanic wraps a middleware that recovers from panics
func NewPanic(opts PanicOptions) webserver.MiddlewareFunc {
	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) (err error) {
			// Define a panic handler
			defer func() {
				if r := recover(); r != nil {
					var ok bool
					var stack []byte
					var stackLen int

					if r == http.ErrAbortHandler {
						panic(r)
					}

					// Get error
					err, ok = r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					// Get stack trace
					if opts.StackSize > 0 {
						stack = make([]byte, opts.StackSize)
						stackLen = runtime.Stack(stack, opts.IncludeAllGoRoutines)
						stack = stack[:stackLen]
					}

					// Call panic error handler
					if opts.PanicErrorHandler != nil {
						err = opts.PanicErrorHandler(req, err, stack)
					} else {
						err = fmt.Errorf("[UNHANDLED EXCEPTION] %v %s\n", err, string(stack))
					}
				}
			}()

			// Run next middleware
			err = next(req)

			// Done
			return
		}
	}
}
