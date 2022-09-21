package schemax

import (
	"log"
	"testing"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

func TestGenerateSchemas(t *testing.T) {
	t.Run("generate article schema from a post", func(t *testing.T) {
		now := time.Now()
		articleSchema := GetArticleSchema(PostData{
			Post: Post{
				Base:          Base{CreatedAt: now},
				Title:         "Test Title",
				PublishedDate: &now,
			},
			Authors: []PostAuthor{
				{
					FirstName: "Test",
					LastName:  "User",
					Slug:      "test-user",
				},
			},
		}, Space{
			Name: "Test Space",
			SpaceSettings: &SpaceSettings{
				SiteTitle: "testspace.com",
				Logo: &Medium{
					URL: postgres.Jsonb{RawMessage: []byte(`{"raw":"testmedium.com/q/1"}`)},
				},
			},
		})

		log.Println(articleSchema)
	})
}
