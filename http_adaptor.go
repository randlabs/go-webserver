package go_webserver

import (
	"net/http"

	"github.com/randlabs/go-webserver/request"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// -----------------------------------------------------------------------------

func HandlerFromHttpHandler(handler http.Handler) HandlerFunc {
	fasthttpHandler := fasthttpadaptor.NewFastHTTPHandler(handler)
	return func(req *request.RequestContext) error {
		req.CallFastHttpHandler(fasthttpHandler)
		return nil
	}
}

func HandlerFromHttpHandlerFunc(f http.HandlerFunc) HandlerFunc {
	return HandlerFromHttpHandler(f)
}
