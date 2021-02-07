package meilisearchx

import "fmt"

// DeleteDocument updates the document in meili index
func DeleteDocument(indexName string, id uint, kind string) error {
	objectID := fmt.Sprint(kind, "_", id)
	_, err := Client.Documents(indexName).Delete(objectID)
	if err != nil {
		return err
	}

	return nil
}
