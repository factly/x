package meilisearchx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/spf13/viper"
)

// SearchWithoutQuery calls meili without q
func SearchWithoutQuery(indexName string, filters string, kind string) (map[string]interface{}, error) {

	body := map[string]interface{}{
		"filters":      filters,
		"facetFilters": []string{"kind:" + kind},
		"limit":        1000000,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", viper.GetString("meili_url")+"/indexes/"+indexName+"/search", buf)
	req.Header.Add("X-Meili-API-Key", viper.GetString("meili_key"))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SearchWithQuery calls meili with q
func SearchWithQuery(indexName, q, filters, kind string) ([]interface{}, error) {
	result, err := Client.Search(indexName).Search(meilisearch.SearchRequest{
		Query:        q,
		Filters:      filters,
		FacetFilters: []string{"kind:" + kind},
		Limit:        1000000,
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
