package meilisearchx

import (
	"time"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/spf13/viper"
)

// Client client for meili search server
var Client *meilisearch.Client
var searchable_attributes []string
var filterable_attributes []string
var sortable_attributes []string
var ranking_attributes []string
var stop_words []string

// SetupMeiliSearch setups the meili search server index
func SetupMeiliSearch(indexes, searchableAttributes, filterableAttributes, sortableAttributes, rankingAttritubes, stopWords []string) error {
	Client = meilisearch.NewClient(
		meilisearch.ClientConfig{
			Host:    viper.GetString("meili_url"),
			APIKey:  viper.GetString("meili_api_key"),
			Timeout: time.Second * 10,
		},
	)

	searchable_attributes = searchableAttributes
	filterable_attributes = filterableAttributes
	sortable_attributes = sortableAttributes
	ranking_attributes = rankingAttritubes
	stop_words = stopWords

	for _, indexName := range indexes {
		setupIndex(indexName)
	}

	return nil
}

func setupIndex(indexName string) error {

	_, err := Client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexName,
		PrimaryKey: "object_id",
	})

	if err != nil {
		return err
	}
	// adding filterable attributes to the meilisearch instance
	_, err = Client.Index(indexName).UpdateFilterableAttributes(&filterable_attributes)
	if err != nil {
		return err
	}

	// Add searchable attributes in index
	_, err = Client.Index(indexName).UpdateSearchableAttributes(&searchable_attributes)
	if err != nil {
		return err
	}

	// Add sortable attributes in index
	_, err = Client.Index(indexName).UpdateSortableAttributes(&sortable_attributes)
	if err != nil {
		return err
	}
	// Add updateTypoTolerance in index
	_, err = Client.Index(indexName).UpdateTypoTolerance(&meilisearch.TypoTolerance{
		MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
			OneTypo:  4,
			TwoTypos: 8,
		},
	})
	if err != nil {
		return err
	}

	// Add UpdateRankingRules in index
	_, err = Client.Index(indexName).UpdateRankingRules(&ranking_attributes)

	if err != nil {
		return err
	}

	// Add sortable attributes in index
	_, err = Client.Index(indexName).UpdateStopWords(&stop_words)
	if err != nil {
		return err
	}
	return nil
}
