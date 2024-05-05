package go_webserver

import (
	httpprof "net/http/pprof"
	"runtime/pprof"
	"strings"
)

// -----------------------------------------------------------------------------

// ServeDebugProfiles adds the GO runtime profile handlers to a web server
func (srv *Server) ServeDebugProfiles(basePath string, middlewares ...HandlerFunc) {
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	srv.GET(basePath, NewHandlerFromHttpHandlerFunc(httpprof.Index), middlewares...)

	for _, profile := range pprof.Profiles() {
		h := httpprof.Handler(profile.Name())
		srv.GET(basePath+profile.Name(), NewHandlerFromHttpHandler(h), middlewares...)
	}
	srv.GET(basePath+"cmdline", NewHandlerFromHttpHandlerFunc(httpprof.Cmdline), middlewares...)
	srv.GET(basePath+"profile", NewHandlerFromHttpHandlerFunc(httpprof.Profile), middlewares...)
	srv.GET(basePath+"symbol", NewHandlerFromHttpHandlerFunc(httpprof.Symbol), middlewares...)
	srv.GET(basePath+"trace", NewHandlerFromHttpHandlerFunc(httpprof.Trace), middlewares...)
}
