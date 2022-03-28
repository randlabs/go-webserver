package request

import "github.com/valyala/fasthttp"

// -----------------------------------------------------------------------------

func (req *RequestContext) reset() {
	req.ctx = nil
}

func (req *RequestContext) sendError(statusCode int, msg string) {
	if !req.IsHead() {
		if len(msg) == 0 {
			msg = fasthttp.StatusMessage(statusCode)
		}
		req.ctx.Error(msg, statusCode)
	} else {
		req.ctx.Error("", statusCode)
	}
}
