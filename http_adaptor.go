package go_webserver

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// -----------------------------------------------------------------------------

func FastHttpHandlerFromHttpHandler(handler http.Handler) fasthttp.RequestHandler {
	// Create handler
	return fasthttpadaptor.NewFastHTTPHandler(handler)
}

func FastHttpHandlerFromHttpHandlerFunc(f http.HandlerFunc) fasthttp.RequestHandler {
	return FastHttpHandlerFromHttpHandler(f)
}
