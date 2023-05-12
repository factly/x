package meilisearchx

import (
	"errors"
	"fmt"
)

// AddDocument addes object into meili search index
func AddDocument(indexName string, data map[string]interface{}) error {
	_, err := Client.GetIndex(indexName)
	if err != nil {
		setupIndex(indexName)
	}
	if data["id"] == nil || data["id"] == "" {
		return errors.New("no id field in meili document")
	}

	data["object_id"] = fmt.Sprint(data["id"])

	arr := []map[string]interface{}{data}

	_, err = Client.Index(indexName).UpdateDocuments(arr)
	if err != nil {
		return err
	}

	return nil
}
