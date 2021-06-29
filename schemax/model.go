package schemax

import (
	"time"

	"github.com/factly/dega-server/service/core/model"
	factCheckModel "github.com/factly/dega-server/service/fact-check/model"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type postData struct {
	model.Post
	Authors []model.Author         `json:"authors"`
	Claims  []factCheckModel.Claim `json:"claims"`
}

// Author of Factcheck
type Author struct {
	Name string `json:"name"`
	Type string `json:"@type"`
	URL  string `json:"url"`
}

// ItemReviewed type
type ItemReviewed struct {
	Type          string         `json:"@type"`
	DatePublished time.Time      `json:"datePublished"`
	Appearance    postgres.Jsonb `json:"appearance"`
	Author        Author         `json:"author"`
}

// ReviewRating type
type ReviewRating struct {
	Type               string `json:"@type"`
	RatingValue        int    `json:"ratingValue"`
	BestRating         int    `json:"bestRating"`
	WorstRating        int    `json:"worstRating"`
	AlternateName      string `json:"alternateName"`
	RatingExplaination string `json:"ratingExplaination"`
}

// FactCheckSchema for factcheck
type FactCheckSchema struct {
	Context       string       `json:"@context"`
	Type          string       `json:"@type"`
	DatePublished time.Time    `json:"datePublished"`
	URL           string       `json:"url"`
	ClaimReviewed string       `json:"claimReviewed"`
	Author        Author       `json:"author"`
	ReviewRating  ReviewRating `json:"reviewRating"`
	ItemReviewed  ItemReviewed `json:"itemReviewed"`
}

// Image for article
type Image struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
}

// Publisher for article
type Publisher struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	Logo Image  `json:"logo"`
}

// ArticleSchema for article
type ArticleSchema struct {
	Context       string     `json:"@context"`
	Type          string     `json:"@type"`
	Headline      string     `json:"headline"`
	Image         []Image    `json:"image"`
	DatePublished *time.Time `json:"datePublished"`
	Author        []Author   `json:"author"`
	Publisher     Publisher  `json:"publisher"`
}
