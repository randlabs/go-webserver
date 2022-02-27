package go_webserver

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

func SendSuccess(ctx *RequestCtx) {
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func SendJSON(ctx *RequestCtx, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	err := json.NewEncoder(ctx).Encode(obj)
	if err != nil {
		SendInternalServerError(ctx, "")
	}
}

func SendBadRequest(ctx *RequestCtx, msg string) {
	sendError(ctx, fasthttp.StatusBadRequest, msg)
}

func SendAccessDenied(ctx *RequestCtx, msg string) {
	sendError(ctx, fasthttp.StatusForbidden, msg)
}

func SendInternalServerError(ctx *RequestCtx, msg string) {
	sendError(ctx, fasthttp.StatusInternalServerError, msg)
}

func EnableCORS(ctx *RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "Content-Type, X-Access-Token, Authorization")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET")
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
}

func DisableCache(ctx *RequestCtx) {
	ctx.Response.Header.Add(
		"Cache-Control",
		"private,no-cache,no-store,max-age=0,must-revalidate,pre-check=0,post-check=0",
	)
}

// -----------------------------------------------------------------------------
// Private functions

func sendError(ctx *RequestCtx, statusCode int, msg string) {
	if len(msg) == 0 {
		msg = fasthttp.StatusMessage(statusCode)
	}
	ctx.Error(msg, statusCode)
}
