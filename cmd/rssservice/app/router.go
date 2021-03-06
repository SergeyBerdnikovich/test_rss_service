package app

import (
	"github.com/SergeyBerdnikovich/test_rss_service/pkg/handlers"
	"github.com/buaazp/fasthttprouter"
)

func initializeRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/authenticate", loggingMiddleware(handlers.NewAuthHandler().Authenticate))
	router.GET("/rss_feeds_items", allMiddlewares(handlers.NewRssHandler().FetchRssFeedsItems))

	return router
}
