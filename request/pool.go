package request

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

type RequestContextPool struct {
	pool sync.Pool
}

// -----------------------------------------------------------------------------

func NewRequestContextPool() *RequestContextPool {
	return &RequestContextPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(RequestContext)
			},
		},
	}
}

func (rcp *RequestContextPool) NewRequestContext(ctx *fasthttp.RequestCtx) (req *RequestContext) {
	req, _ = rcp.pool.Get().(*RequestContext)
	req.ctx = ctx
	return req
}

func (rcp *RequestContextPool) ReleaseRequestContext(req *RequestContext) {
	req.reset()
	rcp.pool.Put(req)
}
