package schemax

import (
	"encoding/json"
	"fmt"
)

func GetArticleSchema(obj PostData, space Space) ArticleSchema {
	jsonLogo := map[string]string{}
	if space.SpaceSettings.Logo != nil {
		rawLogo, _ := space.SpaceSettings.Logo.URL.RawMessage.MarshalJSON()
		_ = json.Unmarshal(rawLogo, &jsonLogo)
	}

	articleSchema := ArticleSchema{}
	articleSchema.Context = "https://schema.org"
	articleSchema.Type = "NewsArticle"
	articleSchema.Headline = obj.Post.Title
	if _, ok := jsonLogo["raw"]; ok {
		articleSchema.Image = append(articleSchema.Image, Image{
			Type: "ImageObject",
			URL:  jsonLogo["raw"]})
	}
	articleSchema.DatePublished = obj.Post.PublishedDate
	for _, eachAuthor := range obj.Authors {
		articleSchema.Author = append(articleSchema.Author, Author{
			Type: "Person",
			Name: eachAuthor.FirstName + " " + eachAuthor.LastName,
			URL:  fmt.Sprint(space.SpaceSettings.SiteAddress, "/users/", eachAuthor.Slug),
		})
	}
	articleSchema.Publisher.Type = "Organization"
	articleSchema.Publisher.Name = space.Name
	if _, ok := jsonLogo["raw"]; ok {
		articleSchema.Publisher.Logo.Type = "ImageObject"
		articleSchema.Publisher.Logo.URL = jsonLogo["raw"]
	}

	return articleSchema
}

func GetFactCheckSchema(obj PostData, space Space, ratings []Rating) []FactCheckSchema {
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
		claimSchema.URL = space.SpaceSettings.SiteAddress + "/" + obj.Slug
		claimSchema.ClaimReviewed = each.Claim
		claimSchema.Author.Type = "Organization"
		claimSchema.Author.Name = space.Name
		claimSchema.Author.URL = space.SpaceSettings.SiteAddress
		claimSchema.ReviewRating.Type = "Rating"
		claimSchema.ReviewRating.RatingValue = each.Rating.NumericValue
		claimSchema.ReviewRating.AlternateName = each.Rating.Name
		claimSchema.ReviewRating.BestRating = bestRating
		claimSchema.ReviewRating.RatingExplanation = each.Fact
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

func GetSchemas(obj PostData, space Space, ratings []Rating) []interface{} {
	schemas := make([]interface{}, 0)

	schemas = append(schemas, GetArticleSchema(obj, space))

	factCheckSchemas := GetFactCheckSchema(obj, space, ratings)

	for _, each := range factCheckSchemas {
		schemas = append(schemas, each)
	}

	return schemas
}
