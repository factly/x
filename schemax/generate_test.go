package schemax

import (
	"log"
	"testing"
	"time"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func TestGenerateSchemas(t *testing.T) {
	t.Run("generate article schema from a post", func(t *testing.T) {
		now := time.Now()
		articleSchema := GetArticleSchema(PostData{
			Post: model.Post{
				Base:          config.Base{CreatedAt: now},
				Title:         "Test Title",
				PublishedDate: &now,
			},
			Authors: []model.Author{
				{
					FirstName: "Test",
					LastName:  "User",
					Slug:      "test-user",
				},
			},
		}, model.Space{
			Name:      "Test Space",
			SiteTitle: "testspace.com",
			Logo: &model.Medium{
				URL: postgres.Jsonb{RawMessage: []byte(`{"raw":"testmedium.com/q/1"}`)},
			},
		})

		log.Println(articleSchema)
	})
}
