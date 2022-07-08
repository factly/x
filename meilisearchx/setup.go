package meilisearchx

import (
	"time"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/spf13/viper"
)

// Client client for meili search server
var Client *meilisearch.Client

// SetupMeiliSearch setups the meili search server index
func SetupMeiliSearch(indexName string, searchableAttributes []string) error {
	Client = meilisearch.NewClient(
		meilisearch.ClientConfig{
			Host:    viper.GetString("meili_url"),
			APIKey:  viper.GetString("meili_api_key"),
			Timeout: time.Second * 10,
		},
	)

	_, err := Client.GetIndex(indexName)
	if err != nil {
		_, err = Client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        indexName,
			PrimaryKey: "object_id",
		})

		if err != nil {
			return err
		}
	}

	// adding filterable attributes to the meilisearch instance
	filterableAttributes := []string{"kind"}
	_, err = Client.Index(indexName).UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		return err
	}

	// Add searchable attributes in index
	_, err = Client.Index(indexName).UpdateSearchableAttributes(&searchableAttributes)
	if err != nil {
		return err
	}
	return nil
}
