package itemsfetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItemsFetcher_FetchItems_SuccessFailAndTimeout(t *testing.T) {
	successResponse1 := `<rss version="2.0">
	<channel>
	<title>Feed Test Title 1</title>
	<link>http://test_feed_link</link>
	<item>
	<title>Test Title 1</title>
	<link>http://test_link1</link>
	<pubDate>Sun, 03 May 2020 22:25:48 PDT</pubDate>
	<description>Test Description 1</description>
	</item>
	<item>
	<title>Test Title 2</title>
	<link>http://test_link2</link>
	<pubDate>Sun, 02 May 2020 11:25:48 PDT</pubDate>
	<description>Test Description 2</description>
	</item>
	</channel>
	</rss>`
	successResponse2 := `<rss version="2.0">
	<channel>
	<title>Feed Test Title 2</title>
	<link>http://test_feed_link</link>
	<item>
	<title>Test Title 3</title>
	<link>http://test_link3</link>
	<pubDate>Sun, 13 May 2020 22:25:48 PDT</pubDate>
	<description>Test Description 3</description>
	</item>
	</channel>
	</rss>`

	server1 := mockServerResponse(200, successResponse1, 2*time.Second)
	server2 := mockServerResponse(404, "", 12*time.Second)
	server3 := mockServerResponse(500, "", 0)
	server4 := mockServerResponse(200, successResponse1, 15*time.Second)
	server5 := mockServerResponse(200, successResponse2, 2*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	urls := []string{
		server1.URL,
		server2.URL,
		server3.URL,
		server4.URL,
		server5.URL,
	}

	items, err := NewItemsFetcher().FetchItems(ctx, urls)

	assert.Equal(t, 3, len(items))
	assert.NotNil(t, err)

	assert.True(t, strings.Contains(items[0].Title, "Test Title"))
	assert.True(t, strings.Contains(items[0].Source, "Feed Test Title"))
	assert.True(t, strings.Contains(items[0].SourceURL, "http://test_feed_link"))
	assert.True(t, strings.Contains(items[0].Link, "http://test_link"))
	assert.True(t, strings.Contains(items[0].Description, "Test Description"))
	assert.True(t, strings.Contains(items[0].PublishDate.String(), "2020-05"))

	assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("%v: Get \"%v\": context deadline exceeded", server2.URL, server2.URL)))
	assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("%v: http error: 500 Internal Server Error", server3.URL)))
	assert.True(t, strings.Contains(err.Error(), fmt.Sprintf("%v: Get \"%v\": context deadline exceeded", server4.URL, server4.URL)))
}

func mockServerResponse(code int, body string, delay time.Duration) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/xml")
		io.WriteString(w, body)
	}))

	return server
}
