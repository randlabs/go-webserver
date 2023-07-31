package middleware

import (
	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// ConditionEvaluator defines a function that executes the wrapped middleware if returns true
type ConditionEvaluator func(req *request.RequestContext) bool

// -----------------------------------------------------------------------------

// NewConditional wraps a middleware to conditionally execute or skip it depending on the evaluator's return value
func NewConditional(cond ConditionEvaluator, m webserver.MiddlewareFunc) webserver.MiddlewareFunc {
	// Setup middleware function
	return func(next webserver.HandlerFunc) webserver.HandlerFunc {
		return func(req *request.RequestContext) error {
			// Evaluate condition
			if cond(req) {
				return m(next)(req)
			}

			// Go to next middleware
			return next(req)
		}
	}
}
