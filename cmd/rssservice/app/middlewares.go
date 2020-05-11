package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SergeyBerdnikovich/test_rss_service/pkg/handlers"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

func allMiddlewares(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return loggingMiddleware(
		authMiddleware(
			handler))
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

func authMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		tokenString := string(ctx.Request.Header.Peek("Token"))
		if tokenString == "" {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return handlers.MyJWTSigningKey, nil
		})
		if err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expiration := claims["exp"]
			if time.Now().Unix() > int64(expiration.(float64)) {
				ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
				return
			}
			handler(ctx)
		} else {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		}
	})
}
