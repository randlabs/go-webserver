package go_webserver

import (
	"net/http"
	httpprof "net/http/pprof"
	"runtime/pprof"
	"strings"

	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

// ProfilerHandlerCheckAccess specifies a callback function that evaluates access to the profiler handlers
type ProfilerHandlerCheckAccess func(ctx *RequestCtx) bool

// -----------------------------------------------------------------------------

// AddProfilerHandlers adds the GO runtime profile handlers to a web server
func (srv *Server) AddProfilerHandlers(basePath string, accessCheck ProfilerHandlerCheckAccess) {
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	srv.Router.GET(basePath, wrapProfilerHandlerFunc(httpprof.Index, accessCheck))

	for _, profile := range pprof.Profiles() {
		h := httpprof.Handler(profile.Name())
		srv.Router.GET(basePath+profile.Name(), wrapProfilerHandler(h, accessCheck))
	}
	srv.Router.GET(basePath+"cmdline", wrapProfilerHandlerFunc(httpprof.Cmdline, accessCheck))
	srv.Router.GET(basePath+"profile", wrapProfilerHandlerFunc(httpprof.Profile, accessCheck))
	srv.Router.GET(basePath+"symbol", wrapProfilerHandlerFunc(httpprof.Symbol, accessCheck))
	srv.Router.GET(basePath+"trace", wrapProfilerHandlerFunc(httpprof.Trace, accessCheck))
}

// -----------------------------------------------------------------------------
// Private functions

func wrapProfilerHandler(handler http.Handler, accessCheck ProfilerHandlerCheckAccess) fasthttp.RequestHandler {
	return wrapProfilerFastHandler(FastHttpHandlerFromHttpHandler(handler), accessCheck)
}

func wrapProfilerHandlerFunc(handler http.HandlerFunc, accessCheck ProfilerHandlerCheckAccess) fasthttp.RequestHandler {
	return wrapProfilerFastHandler(FastHttpHandlerFromHttpHandlerFunc(handler), accessCheck)
}

func wrapProfilerFastHandler(handler fasthttp.RequestHandler, accessCheck ProfilerHandlerCheckAccess) fasthttp.RequestHandler {
	wrapper := func(ctx *RequestCtx) {
		// Disable cache for this requests
		EnableCORS(ctx)
		DisableCache(ctx)

		// Check access
		if accessCheck != nil && (!accessCheck(ctx)) {
			// Deny access
			SendAccessDenied(ctx, "403 forbidden")
			return
		}

		// Call the handler
		handler(ctx)
	}
	return wrapper
}
