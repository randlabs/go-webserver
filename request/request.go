package request

//go:generate go run generators/request_inherited.go

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

type RequestContext struct {
	ctx *fasthttp.RequestCtx
}

// -----------------------------------------------------------------------------

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

// -----------------------------------------------------------------------------

func (req *RequestContext) Request() *fasthttp.Request {
	return &req.ctx.Request
}

func (req *RequestContext) Response() *fasthttp.Response {
	return &req.ctx.Response
}

func (req *RequestContext) RequestHeaders() *fasthttp.RequestHeader {
	return &req.ctx.Request.Header
}

func (req *RequestContext) ResponseHeaders() *fasthttp.ResponseHeader {
	return &req.ctx.Response.Header
}

func (req *RequestContext) RequestHeader(key string) string {
	return string(req.ctx.Request.Header.Peek(key))
}

func (req *RequestContext) SetResponseHeader(key string, value string) {
	req.ctx.Response.Header.Set(key, value)
}

func (req *RequestContext) AddResponseHeader(key string, value string) {
	req.ctx.Response.Header.Add(key, value)
}

func (req *RequestContext) SendSuccess() {
	req.ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func (req *RequestContext) NotFound(msg string) {
	req.sendError(fasthttp.StatusNotFound, msg)
}

func (req *RequestContext) BadRequest(msg string) {
	req.sendError(fasthttp.StatusBadRequest, msg)
}

func (req *RequestContext) AccessDenied(msg string) {
	req.sendError(fasthttp.StatusForbidden, msg)
}

func (req *RequestContext) InternalServerError(msg string) {
	req.sendError(fasthttp.StatusInternalServerError, msg)
}

func (req *RequestContext) SendJSON(obj interface{}) {
	req.ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	req.ctx.Response.SetStatusCode(fasthttp.StatusOK)

	err := json.NewEncoder(req.ctx).Encode(obj)
	if err != nil {
		req.InternalServerError("unable to encode json output")
	}
}

func (req *RequestContext) NoContent(statusCode int) error {
	req.ctx.Response.SetStatusCode(statusCode)
	req.ctx.Response.ResetBody()
	return nil
}

func (req *RequestContext) CallFastHttpHandler(handler fasthttp.RequestHandler) {
	handler(req.ctx)
}
