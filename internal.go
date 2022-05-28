package go_webserver

import (
	"net/http"

	"github.com/randlabs/go-webserver/request"
	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

const (
	requestContextUV = "\xFF\xFF**request"
	requestErrorUV   = "\xFF\xFF**reqError"
)

// -----------------------------------------------------------------------------

func (srv *Server) createMasterHandler(masterHandler fasthttp.RequestHandler) fasthttp.RequestHandler {
	// Create wrapper for the master handler
	masterHandlerWrapper := func(req *request.RequestContext) error {
		req.CallFastHttpHandler(masterHandler)

		err, ok := req.UserValue(requestErrorUV).(error)
		if ok {
			req.RemoveUserValue(requestErrorUV)
			return err
		}
		return nil
	}

	// Build the recursive function for server-wide middlewares
	var f func(idx int) HandlerFunc
	f = func(idx int) HandlerFunc {
		if idx < len(srv.middlewares) {
			return srv.middlewares[idx](f(idx + 1))
		}
		return masterHandlerWrapper
	}

	// Wrapper
	return func(ctx *fasthttp.RequestCtx) {
		// Create a new request context
		req := srv.requestPool.NewRequestContext(ctx)

		// Bind request object to fasthttp.RequestCtx
		ctx.SetUserValue(requestContextUV, req)

		// Get the top handler in the chain
		topHandler := f(0)

		// Call it
		err := topHandler(req)
		if err != nil {
			srv.requestErrorHandler(req, err)
		}

		// Unbind request object from fasthttp.RequestCtx
		ctx.RemoveUserValue(requestContextUV)

		// Release request context
		srv.requestPool.ReleaseRequestContext(req)
	}
}

func (srv *Server) createEndpointHandler(epHandler HandlerFunc, middlewares ...MiddlewareFunc) fasthttp.RequestHandler {
	// Build the recursive function for endpoint middlewares
	var f func(idx int) HandlerFunc

	if middlewares == nil || len(middlewares) == 0 {
		f = func(_ int) HandlerFunc {
			return epHandler
		}
	} else {
		f = func(idx int) HandlerFunc {
			if idx < len(middlewares) {
				return middlewares[idx](f(idx + 1))
			}
			return epHandler
		}
	}

	// Wrapper
	return func(ctx *fasthttp.RequestCtx) {
		// Obtain the attached RequestContext object
		req, ok := ctx.UserValue(requestContextUV).(*request.RequestContext)
		if !ok {
			// This should not happen but just in case the RequestContext object is not found
			ctx.Error("unhandled endpoint", fasthttp.StatusInternalServerError)
			return
		}

		// Get the top handler in the chain
		topHandler := f(0)

		// Call it
		err := topHandler(req)

		// On error, save it so the master handler can process it
		if err != nil {
			ctx.SetUserValue(requestErrorUV, err)
		}
	}
}

func (srv *Server) defaultRequestErrorHandler(req *request.RequestContext, err error) {
	req.SetStatusCode(http.StatusInternalServerError)
	req.ResetBody()
	req.InternalServerError(err.Error())
}

func (srv *Server) logCallback(format string, args ...interface{}) {
	// Nothing for now
}
