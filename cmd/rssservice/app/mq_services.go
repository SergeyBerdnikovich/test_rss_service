package app

import (
	"encoding/json"
	"time"

	"github.com/SergeyBerdnikovich/test_rss_service/pkg/entities"
	"github.com/SergeyBerdnikovich/test_rss_service/pkg/services/mqpublisher"
	"github.com/SergeyBerdnikovich/test_rss_service/pkg/services/rss/itemsfetcher"
	"golang.org/x/net/context"
)

func initializeMQServices() {
	mqpublisher.MQPublisher.InitializeQueues()

	rssFeedsService()
}

func closeMQServices() func() {
	return func() {
		mqpublisher.MQPublisher.Close()
	}
}

func rssFeedsService() {
	msgs := mqpublisher.MQPublisher.Consume("rss.feed_urls")

	go func() {
		for d := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
			defer cancel()

			var feedURL string
			json.Unmarshal(d.Body, &feedURL)

			items, err := itemsfetcher.NewItemsFetcher().FetchItems(ctx, []string{feedURL})

			type rssItemsResponse struct {
				Items []entities.RssItem `json:"items"`
				Err   string             `json:"errors,omitempty"`
			}

			data, _ := json.Marshal(rssItemsResponse{Items: items, Err: parseErrorMessages(err)})

			mqpublisher.MQPublisher.Publish("rss.feed_items", data)
		}
	}()
}

func parseErrorMessages(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
