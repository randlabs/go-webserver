package go_webserver

import (
	"github.com/randlabs/go-webserver/middleware"
	"net/http"
	httpprof "net/http/pprof"
	"runtime/pprof"
	"strings"

	"github.com/randlabs/go-webserver/models"
	"github.com/randlabs/go-webserver/request"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// -----------------------------------------------------------------------------

// ProfilerHandlerCheckAccess specifies a callback function that evaluates access to the profiler handlers
type ProfilerHandlerCheckAccess func(req *request.RequestContext) bool

// -----------------------------------------------------------------------------

// AddProfilerHandlers adds the GO runtime profile handlers to a web server
func (srv *Server) AddProfilerHandlers(
	basePath string, accessCheck ProfilerHandlerCheckAccess, middlewares ...middleware.MiddlewareFunc,
) {
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	srv.GET(basePath, wrapProfilerHandlerFunc(httpprof.Index, accessCheck), middlewares...)

	for _, profile := range pprof.Profiles() {
		h := httpprof.Handler(profile.Name())
		srv.GET(basePath+profile.Name(), wrapProfilerHandler(h, accessCheck), middlewares...)
	}
	srv.GET(basePath+"cmdline", wrapProfilerHandlerFunc(httpprof.Cmdline, accessCheck), middlewares...)
	srv.GET(basePath+"profile", wrapProfilerHandlerFunc(httpprof.Profile, accessCheck), middlewares...)
	srv.GET(basePath+"symbol", wrapProfilerHandlerFunc(httpprof.Symbol, accessCheck), middlewares...)
	srv.GET(basePath+"trace", wrapProfilerHandlerFunc(httpprof.Trace, accessCheck), middlewares...)
}

// -----------------------------------------------------------------------------
// Private functions

func wrapProfilerHandler(handler http.Handler, accessCheck ProfilerHandlerCheckAccess) models.HandlerFunc {
	fasthttpHandler := fasthttpadaptor.NewFastHTTPHandler(handler)

	return func(req *request.RequestContext) error {
		// Disable cache for this requests
		//		EnableCORS(ctx)
		//		DisableCache(ctx)

		// Check access
		if accessCheck == nil || accessCheck(req) {
			// Call the handler
			req.CallFastHttpHandler(fasthttpHandler)
		} else {
			// Deny access
			req.AccessDenied("403 forbidden")
		}

		// Done
		return nil
	}
}

func wrapProfilerHandlerFunc(handler http.HandlerFunc, accessCheck ProfilerHandlerCheckAccess) models.HandlerFunc {
	return wrapProfilerHandler(handler, accessCheck)
}
