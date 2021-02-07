package meilisearchx

import (
	"log"
	"net/http"
	"time"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/spf13/viper"
)

// Client client for meili search server
var Client *meilisearch.Client

// SetupMeiliSearch setups the meili search server index
func SetupMeiliSearch(indexName string, searchableAttributes []string) {
	Client = meilisearch.NewClientWithCustomHTTPClient(meilisearch.Config{
		Host:   viper.GetString("meili_url"),
		APIKey: viper.GetString("meili_key"),
	}, http.Client{
		Timeout: time.Second * 10,
	})

	_, err := Client.Indexes().Get(indexName)
	if err != nil {
		_, err = Client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID:        indexName,
			PrimaryKey: "object_id",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = Client.Settings(indexName).UpdateAttributesForFaceting([]string{"kind"})
	if err != nil {
		log.Fatal(err)
	}

	// Add searchable attributes in index
	_, err = Client.Settings(indexName).UpdateSearchableAttributes(searchableAttributes)
	if err != nil {
		log.Fatal(err)
	}
}

// []string{"name", "slug", "description", "title", "subtitle", "excerpt", "site_title", "site_address", "tag_line", "review", "review_tag_line"}
