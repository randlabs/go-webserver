package models

import (
	"github.com/randlabs/go-webserver/request"
)

// -----------------------------------------------------------------------------

// HandlerFunc defines a function that handles a request
type HandlerFunc func(req *request.RequestContext) error
