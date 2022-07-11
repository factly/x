package meilisearchx

import (
	"fmt"
	"strings"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// SearchWithQuery calls meili with q
func SearchWithQuery(indexName, q, filters, kind string) ([]interface{}, error) {
	filter := [][]string{}
	filter = append(filter, []string{filters})
	if kind != "" {
		filter = append(filter, []string{fmt.Sprintf("kind=%s", kind)})
	}
	result, err := Client.Index(indexName).Search(q, &meilisearch.SearchRequest{
		Filter: filter,
		Limit:  1000000,
	})

	if err != nil {
		return nil, err
	}
	return result.Hits, nil
}

// GetIDArray gets array of IDs for search results
func GetIDArray(hits []interface{}) []uint {
	arr := make([]uint, 0)

	if len(hits) == 0 {
		return arr
	}

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		id := hitMap["id"].(float64)
		arr = append(arr, uint(id))
	}

	return arr
}

// GenerateFieldFilter generates filter in form "(field=x OR field=y OR ...)"
func GenerateFieldFilter(ids []string, field string) string {
	filter := "("
	for i, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" {
			if i == len(ids)-1 {
				filter = fmt.Sprint(filter, field, "=", id, ")")
				break
			}
			filter = fmt.Sprint(filter, field, "=", id, " OR ")
		}
	}
	return filter
}
