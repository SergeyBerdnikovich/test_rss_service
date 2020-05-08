package itemsfetcher

import (
	"context"

	"github.com/SergeyBerdnikovich/test_rss_reader/pkg/rssreader"
	"github.com/SergeyBerdnikovich/test_rss_service/pkg/entities"
	"github.com/pkg/errors"
)

// ItemsFetcher is an interface of itemsFetcher service object
// mockgen -source=./pkg/services/rss/itemsfetcher/items_fetcher.go -destination=./pkg/services/rss/itemsfetcher/items_fetcher_mock.go -package=itemsfetcher itemsFetcher
type ItemsFetcher interface {
	FetchItems(ctx context.Context, urls []string) ([]entities.RssItem, error)
}

type itemsFetcher struct{}

// NewItemsFetcher is a constructor of itemsFetcher service object
func NewItemsFetcher() *itemsFetcher {
	return &itemsFetcher{}
}

func (f *itemsFetcher) FetchItems(ctx context.Context, urls []string) ([]entities.RssItem, error) {
	feedsItems, err := rssreader.Parse(ctx, urls)
	if err != nil {
		err = errors.Wrap(err, "some of feeds are not available")
	}

	items := []entities.RssItem{}
	for _, item := range feedsItems {
		items = append(items, entities.RssItem{
			Title:       item.Title,
			Source:      item.Source,
			SourceURL:   item.SourceURL,
			Link:        item.Link,
			PublishDate: item.PublishDate,
			Description: item.Description,
		})
	}

	return items, err
}
