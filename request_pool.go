package go_webserver

import (
	"sync"

	"github.com/randlabs/go-webserver/v2/trusted_proxy"
	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

type RequestContextPool struct {
	pool sync.Pool
}

// -----------------------------------------------------------------------------

func newRequestContextPool() *RequestContextPool {
	return &RequestContextPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(RequestContext)
			},
		},
	}
}

func (rcp *RequestContextPool) newRequestContext(ctx *fasthttp.RequestCtx, tp *trusted_proxy.TrustedProxy, h HandlerFunc, srvMiddlewares, middlewares []HandlerFunc) (*RequestContext, func()) {
	req, _ := rcp.pool.Get().(*RequestContext)
	req.ctx = ctx
	req.tp = tp

	req.middlewareIndex = 0
	req.srvMiddlewares = srvMiddlewares
	req.middlewares = middlewares
	req.srvMiddlewaresLen = len(srvMiddlewares)
	req.totalMiddlewaresLen = req.srvMiddlewaresLen + len(middlewares)
	req.handler = h

	return req, func() {
		req.ctx = nil
		req.tp = nil
		req.userCtx = nil
		// req.middlewareIndex = 0
		req.srvMiddlewares = nil
		req.middlewares = nil
		// req.srvMiddlewaresLen = 0
		// req.totalMiddlewaresLen = 0
		req.handler = nil

		rcp.pool.Put(req)
	}
}
