package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/SergeyBerdnikovich/test_rss_service/pkg/entities"
	"github.com/SergeyBerdnikovich/test_rss_service/pkg/services/rss/itemsfetcher"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const fetchFeedsTimeout = 6 * time.Second

type rssHandler struct {
	itemsFetcher itemsfetcher.ItemsFetcher
}

// NewRssHandler is a constructor for rssHandler
func NewRssHandler() *rssHandler {
	return &rssHandler{
		itemsFetcher: itemsfetcher.NewItemsFetcher(),
	}
}

// FetchRssFeedsItems is a function which handles get rss items request
// Route for it is GET /rss_feeds_items
// Incoming params is urls which is slice of strings
// The response is json of fetched rss items and the errors if they are exists
func (h *rssHandler) FetchRssFeedsItems(ctx *fasthttp.RequestCtx) {
	type rssItemsResponse struct {
		Items []entities.RssItem `json:"items"`
		Err   string             `json:"errors,omitempty"`
	}

	urls := []string{}
	dataBytes := ctx.QueryArgs().PeekMulti("urls[]")
	for _, data := range dataBytes {
		urls = append(urls, string(data))
	}

	hctx, cancel := context.WithTimeout(ctx, fetchFeedsTimeout)
	defer cancel()

	rssItems, err := h.itemsFetcher.FetchItems(hctx, urls)
	if err != nil {
		logrus.Infof("Failed to fetch some feeds, err: %v", err)
	}

	respondWithJSON(ctx, 200, rssItemsResponse{
		Items: rssItems,
		Err:   parseErrorMessages(err),
	})
}

func respondWithJSON(ctx *fasthttp.RequestCtx, code int, data interface{}) {
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.SetStatusCode(200)

	if err := json.NewEncoder(ctx).Encode(data); err != nil {
		logrus.Errorf("Failed to build json, err: %v", err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func parseErrorMessages(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
