package go_webserver

import (
	"net/http"

	"github.com/randlabs/go-webserver/request"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// -----------------------------------------------------------------------------

// HandlerFromHttpHandler returns a HandlerFunc based on the provided http.Handler
func HandlerFromHttpHandler(handler http.Handler) HandlerFunc {
	h := fasthttpadaptor.NewFastHTTPHandler(handler)
	return func(req *request.RequestContext) error {
		req.CallFastHttpHandler(h)
		return nil
	}
}

// HandlerFromHttpHandlerFunc returns a HandlerFunc based on the provided http.HandlerFunc
func HandlerFromHttpHandlerFunc(f http.HandlerFunc) HandlerFunc {
	return HandlerFromHttpHandler(f)
}
