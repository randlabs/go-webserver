package go_webserver

import (
	httpprof "net/http/pprof"
	"runtime/pprof"
	"strings"
)

// -----------------------------------------------------------------------------

// ServeDebugProfiles adds the GO runtime profile handlers to a web server
func (srv *Server) ServeDebugProfiles(basePath string, middlewares ...MiddlewareFunc) {
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	srv.GET(basePath, HandlerFromHttpHandlerFunc(httpprof.Index), middlewares...)

	for _, profile := range pprof.Profiles() {
		h := httpprof.Handler(profile.Name())
		srv.GET(basePath+profile.Name(), HandlerFromHttpHandler(h), middlewares...)
	}
	srv.GET(basePath+"cmdline", HandlerFromHttpHandlerFunc(httpprof.Cmdline), middlewares...)
	srv.GET(basePath+"profile", HandlerFromHttpHandlerFunc(httpprof.Profile), middlewares...)
	srv.GET(basePath+"symbol", HandlerFromHttpHandlerFunc(httpprof.Symbol), middlewares...)
	srv.GET(basePath+"trace", HandlerFromHttpHandlerFunc(httpprof.Trace), middlewares...)
}
