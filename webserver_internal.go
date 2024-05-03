package go_webserver

import (
	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

func (srv *Server) createEndpointHandler(h HandlerFunc, middlewares ...HandlerFunc) fasthttp.RequestHandler {
	// Wrapper
	return func(ctx *fasthttp.RequestCtx) {
		// Create a new request context
		req, freeReq := srv.requestCtxPool.newRequestContext(ctx, srv.trustedProxy, h, srv.middlewares, middlewares)
		defer freeReq()

		err := req.Next()
		if err != nil {
			srv.requestErrorHandler(req, err)
		}
	}
}
