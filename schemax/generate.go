package schemax

import (
	"encoding/json"
	"fmt"

	"github.com/factly/dega-server/service/core/model"
	factCheckModel "github.com/factly/dega-server/service/fact-check/model"
)

func GetArticleSchema(obj postData, space model.Space) ArticleSchema {
	jsonLogo := map[string]string{}
	if space.Logo != nil {
		rawLogo, _ := space.Logo.URL.RawMessage.MarshalJSON()
		_ = json.Unmarshal(rawLogo, &jsonLogo)
	}

	articleSchema := ArticleSchema{}
	articleSchema.Context = "https://schema.org"
	articleSchema.Type = "NewsArticle"
	articleSchema.Headline = obj.Post.Title
	articleSchema.Image = append(articleSchema.Image, Image{
		Type: "ImageObject",
		URL:  jsonLogo["raw"]})
	articleSchema.DatePublished = obj.Post.PublishedDate
	for _, eachAuthor := range obj.Authors {
		articleSchema.Author = append(articleSchema.Author, Author{
			Type: "Person",
			Name: eachAuthor.FirstName + " " + eachAuthor.LastName,
			URL:  fmt.Sprint(space.SiteAddress, "/users/", eachAuthor.Slug),
		})
	}
	articleSchema.Publisher.Type = "Organization"
	articleSchema.Publisher.Name = space.Name
	articleSchema.Publisher.Logo.Type = "ImageObject"
	articleSchema.Publisher.Logo.URL = jsonLogo["raw"]

	return articleSchema
}

func GetFactCheckSchema(obj postData, space model.Space, ratings []factCheckModel.Rating) []FactCheckSchema {
	result := make([]FactCheckSchema, 0)

	bestRating := 5
	worstRating := 1
	if len(ratings) > 2 {
		bestRating = ratings[len(ratings)-1].NumericValue
		worstRating = ratings[0].NumericValue
	}

	for _, each := range obj.Claims {
		claimSchema := FactCheckSchema{}
		claimSchema.Context = "https://schema.org"
		claimSchema.Type = "ClaimReview"
		claimSchema.DatePublished = obj.Post.CreatedAt
		claimSchema.URL = space.SiteAddress + "/" + obj.Slug
		claimSchema.ClaimReviewed = each.Claim
		claimSchema.Author.Type = "Organization"
		claimSchema.Author.Name = space.Name
		claimSchema.Author.URL = space.SiteAddress
		claimSchema.ReviewRating.Type = "Rating"
		claimSchema.ReviewRating.RatingValue = each.Rating.NumericValue
		claimSchema.ReviewRating.AlternateName = each.Rating.Name
		claimSchema.ReviewRating.BestRating = bestRating
		claimSchema.ReviewRating.RatingExplaination = each.Fact
		claimSchema.ReviewRating.WorstRating = worstRating
		claimSchema.ItemReviewed.Type = "Claim"
		if each.CheckedDate != nil {
			claimSchema.ItemReviewed.DatePublished = *each.CheckedDate
		}
		claimSchema.ItemReviewed.Appearance = each.ClaimSources
		claimSchema.ItemReviewed.Author.Type = "Organization"
		claimSchema.ItemReviewed.Author.Name = each.Claimant.Name

		result = append(result, claimSchema)
	}
	return result
}

func GetSchemas(obj postData, space model.Space, ratings []factCheckModel.Rating) []interface{} {
	schemas := make([]interface{}, 0)

	schemas = append(schemas, GetArticleSchema(obj, space))

	factCheckSchemas := GetFactCheckSchema(obj, space, ratings)

	for _, each := range factCheckSchemas {
		schemas = append(schemas, each)
	}

	return schemas
}
