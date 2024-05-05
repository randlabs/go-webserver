package middleware

import (
	webserver "github.com/randlabs/go-webserver/v2"
)

// -----------------------------------------------------------------------------

// NewNoOP creates a no-operation middleware
func NewNoOP() webserver.HandlerFunc {
	// Setup middleware function
	return func(req *webserver.RequestContext) error {
		// Go to next middleware
		return req.Next()
	}
}
