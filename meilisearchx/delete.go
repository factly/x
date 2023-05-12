package meilisearchx

import "fmt"

// DeleteDocument updates the document in meili index
func DeleteDocument(indexName string, id uint) error {
	objectID := fmt.Sprint(id)
	_, err := Client.Index(indexName).Delete(objectID)
	if err != nil {
		return err
	}

	return nil
}
