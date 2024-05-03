package middleware

import (
	webserver "github.com/randlabs/go-webserver/v2"
)

// -----------------------------------------------------------------------------

// ConditionEvaluator defines a function that executes the wrapped middleware if returns true
type ConditionEvaluator func(req *webserver.RequestContext) (bool, error)

// -----------------------------------------------------------------------------

// NewConditional wraps a middleware to conditionally execute or skip it depending on the evaluator's return value
func NewConditional(cond ConditionEvaluator, m webserver.HandlerFunc) webserver.HandlerFunc {
	// Setup middleware function
	return func(req *webserver.RequestContext) error {
		// Evaluate condition
		ok, err := cond(req)
		if err != nil {
			return err
		}
		if ok {
			return m(req)
		}

		// Go to next middleware
		return req.Next()
	}
}
