package handlers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/SergeyBerdnikovich/test_rss_service/pkg/entities"
	"github.com/SergeyBerdnikovich/test_rss_service/pkg/services/rss/itemsfetcher"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestRssHandler_FetchRssFeedsItems_NoUrls(t *testing.T) {
	t.Run("It returns JSON with empty slice of items", func(t *testing.T) {
		request, err := http.NewRequest("GET", "http://rss_feeds_items", nil)
		if err != nil {
			t.Error(err)
		}

		handler := NewRssHandler()

		response, err := serveHTTPRequest(handler.FetchRssFeedsItems, request)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 200, response.StatusCode)
		assert.True(t, strings.Contains(string(body), "{\"items\":[]}"))
	})
}

func TestRssHandler_FetchRssFeedsItems_ValidAndInvalidUrls(t *testing.T) {
	t.Run("It returns JSON with items and errors", func(t *testing.T) {
		request, err := http.NewRequest("GET", "http://rss_feeds_items?urls=http://valid_url&urls=http://invalid_url", nil)
		if err != nil {
			t.Error(err)
		}

		ctrl := gomock.NewController(t)

		mockItemsFetcher := itemsfetcher.NewMockItemsFetcher(ctrl)

		handler := &rssHandler{
			itemsFetcher: mockItemsFetcher,
		}

		fetchError := errors.New("Test error message")
		rssItems := []entities.RssItem{
			{
				Title:       "Test titile",
				Source:      "Test source",
				SourceURL:   "test_source_url",
				Link:        "test_link",
				Description: "Test description",
			},
		}

		mockItemsFetcher.EXPECT().FetchItems(gomock.Any(), []string{"http://valid_url", "http://invalid_url"}).Return(rssItems, fetchError)

		response, err := serveHTTPRequest(handler.FetchRssFeedsItems, request)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}

		resultInJSON := "{\"items\":[{\"title\":\"Test titile\",\"source\":\"Test source\",\"source_url\":\"test_source_url\",\"link\":\"test_link\",\"publish_date\":null,\"description\":\"Test description\"}],\"errors\":\"Test error message\"}\n"

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, resultInJSON, string(body))
	})
}

func serveHTTPRequest(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}
