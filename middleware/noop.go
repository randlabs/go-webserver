// See the LICENSE file for license details.

package middleware

import (
	webserver "github.com/mxmauro/go-webserver/v2"
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
