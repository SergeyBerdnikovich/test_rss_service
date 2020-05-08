package app

import (
	"encoding/json"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

func allMiddlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return loggingMiddleware(handler)
}

func loggingMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	type logEntry struct {
		Method     string
		RequestURI string
		Status     int
		Host       string
		RemoteAddr string
		UserAgent  string
		Duration   string
	}

	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		handler(ctx)

		duration := time.Since(start)
		reqLog := logEntry{
			Host:       string(ctx.Request.Host()),
			RemoteAddr: ctx.RemoteAddr().String(),
			Method:     string(ctx.Request.Header.Method()),
			RequestURI: string(ctx.Request.URI().RequestURI()),
			UserAgent:  string(ctx.Request.Header.UserAgent()),
			Duration:   duration.String(),
			Status:     ctx.Response.StatusCode(),
		}

		output, _ := json.Marshal(reqLog)
		log.Println(string(output))
	})
}
