package middleware

import (
	"github.com/randlabs/go-webserver/models"
)

// -----------------------------------------------------------------------------

// MiddlewareFunc defines a function that is executed when a request is received
type MiddlewareFunc func(next models.HandlerFunc) models.HandlerFunc
