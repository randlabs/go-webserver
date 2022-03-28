package middleware

import (
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// SkipMiddleware defines a function that indicates if the middleware should be skipped
type SkipMiddleware func(req *request.RequestContext) bool

// -----------------------------------------------------------------------------

func defaultSkip(_ *request.RequestContext) bool {
	return false
}
