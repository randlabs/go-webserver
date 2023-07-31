package middleware

import (
	"crypto/subtle"
	"strings"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// ProtectedEndpointEvaluator evaluates if endpoint access must be denied. Return true to deny access.
type ProtectedEndpointEvaluator func(req *request.RequestContext) bool

// -----------------------------------------------------------------------------

// ProtectedWithToken creates a protection middleware based on an access token string
func ProtectedWithToken(accessToken string) webserver.MiddlewareFunc {
	// Allow access if no token is provided
	if len(accessToken) == 0 {
		return NewNoOP()
	}
	tokenBytes := []byte(accessToken)

	// Create a new protected
	return NewProtected(func(req *request.RequestContext) bool {
		var token []byte

		// Get X-Access-Token header
		header := req.RequestHeader("X-Access-Token")
		if len(header) > 0 {
			token = []byte(header)
		} else {
			// If no token provided, try with Authorization: Bearer XXX header
			header = req.RequestHeader("Authorization")
			if len(header) > 0 {
				auth := strings.SplitN(header, " ", 2)
				if len(auth) == 2 && strings.EqualFold("Bearer", auth[0]) {
					token = []byte(auth[1])
				}
			}
		}

		//Check token
		if len(token) > 0 && subtle.ConstantTimeCompare(tokenBytes, token) != 0 {
			return false // Allow access
		}

		// Deny access
		return true
	})
}

// NewProtected creates a protection middleware based on an evaluator callback
func NewProtected(evaluator ProtectedEndpointEvaluator) webserver.MiddlewareFunc {
	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			// Run evaluator and, if it returns true, assume the endpoint is protected
			if evaluator != nil && evaluator(req) {
				req.AccessDenied("403 forbidden")
				return nil
			}

			// Go to next middleware
			return next(req)
		}
	}
}
