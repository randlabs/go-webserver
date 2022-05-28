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
	tokenBytes := []byte(accessToken)
	return NewProtected(func(req *request.RequestContext) bool {
		if len(tokenBytes) == 0 {
			return false // Allow access
		}

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
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			if evaluator != nil && evaluator(req) {
				req.AccessDenied("403 forbidden")
				return nil
			}
			return next(req)
		}
	}
}
