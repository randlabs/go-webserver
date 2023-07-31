package middleware

import (
	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// NewNoOP creates a no-operation middleware
func NewNoOP() webserver.MiddlewareFunc {
	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			// Go to next middleware
			return next(req)
		}
	}
}
