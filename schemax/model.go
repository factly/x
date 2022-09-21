package schemax

import (
	"time"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

type Base struct {
	ID          uint            `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `sql:"index" json:"deleted_at" swaggertype:"primitive,string"`
	CreatedByID uint            `gorm:"column:created_by_id" json:"created_by_id"`
	UpdatedByID uint            `gorm:"column:updated_by_id" json:"updated_by_id"`
}

type PostData struct {
	Post
	Authors []PostAuthor `json:"authors"`
	Claims  []Claim  `json:"claims"`
}

//Post Author 
type PostAuthor struct {
	Base
	Email            string         `gorm:"column:email;uniqueIndex" json:"email"`
	KID              string         `gorm:"column:kid;" json:"kid"`
	FirstName        string         `gorm:"column:first_name" json:"first_name"`
	LastName         string         `gorm:"column:last_name" json:"last_name"`
	Slug             string         `gorm:"column:slug" json:"slug"`
	DisplayName      string         `gorm:"column:display_name" json:"display_name"`
	BirthDate        string         `gorm:"column:birth_date" json:"birth_date"`
	Gender           string         `gorm:"column:gender" json:"gender"`
	FeaturedMediumID *uint          `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	Medium           *Medium        `gorm:"foreignKey:featured_medium_id" json:"medium"`
	SocialMediaURLs  postgres.Jsonb `gorm:"column:social_media_urls" json:"social_media_urls" swaggertype:"primitive,string"`
	Description      string         `gorm:"column:description" json:"description"`
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
	Type              string `json:"@type"`
	RatingValue       int    `json:"ratingValue"`
	BestRating        int    `json:"bestRating"`
	WorstRating       int    `json:"worstRating"`
	AlternateName     string `json:"alternateName"`
	RatingExplanation string `json:"ratingExplanation"`
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

type Rating struct {
	Base
	Name             string         `gorm:"column:name" json:"name"`
	Slug             string         `gorm:"column:slug" json:"slug"`
	BackgroundColour postgres.Jsonb `gorm:"column:background_colour" json:"background_colour" swaggertype:"primitive,string"`
	TextColour       postgres.Jsonb `gorm:"column:text_colour" json:"text_colour" swaggertype:"primitive,string"`
	Description      postgres.Jsonb `gorm:"column:description" json:"description" swaggertype:"primitive,string"`
	DescriptionHTML  string         `gorm:"column:description_html" json:"description_html,omitempty"`
	NumericValue     int            `gorm:"column:numeric_value" json:"numeric_value"`
	MediumID         *uint          `gorm:"column:medium_id;default=NULL" json:"medium_id"`
	Medium           Medium         `json:"medium"`
	MetaFields       postgres.Jsonb `gorm:"column:meta_fields" json:"meta_fields" swaggertype:"primitive,string"`
	SpaceID          uint           `gorm:"column:space_id" json:"space_id"`
	Meta             postgres.Jsonb `gorm:"column:meta" json:"meta" swaggertype:"primitive,string"`
	HeaderCode       string         `gorm:"column:header_code" json:"header_code"`
	FooterCode       string         `gorm:"column:footer_code" json:"footer_code"`
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
	Logo Image  `json:"logo,omitempty"`
}

// ArticleSchema for article
type ArticleSchema struct {
	Context       string     `json:"@context"`
	Type          string     `json:"@type"`
	Headline      string     `json:"headline"`
	Image         []Image    `json:"image,omitempty"`
	DatePublished *time.Time `json:"datePublished"`
	Author        []Author   `json:"author"`
	Publisher     Publisher  `json:"publisher"`
}

type Space struct {
	Base
	Name            string         `gorm:"column:name" json:"name"`
	Slug            string         `gorm:"column:slug" json:"slug"`
	Description     string         `gorm:"column:description" json:"description"`
	MetaFields      postgres.Jsonb `gorm:"column:meta_fields" json:"meta_fields" swaggertype:"primitive,string"`
	SpaceSettingsID uint           `gorm:"column:space_settings_id;default:NULL" json:"space_settings_id"`
	SpaceSettings   *SpaceSettings `gorm:"foreignKey:space_settings_id" json:"space_settings"`
	OrganisationID  int            `gorm:"column:organisation_id" json:"organisation_id"`
	ApplicationID   uint           `gorm:"column:application_id" json:"application_id"`
}

type SpaceSettings struct {
	Base
	SpaceID           uint           `gorm:"space_id" json:"space_id"`
	SiteTitle         string         `gorm:"site_title" json:"site_title"`
	TagLine           string         `gorm:"tag_line" json:"tag_line"`
	SiteAddress       string         `gorm:"site_address" json:"site_address"`
	LogoID            *uint          `gorm:"column:logo_id;default:NULL" json:"logo_id"`
	Logo              *Medium        `gorm:"foreignKey:logo_id" json:"logo"`
	LogoMobileID      *uint          `gorm:"column:logo_mobile_id;default:NULL" json:"logo_mobile_id"`
	LogoMobile        *Medium        `gorm:"foreignKey:logo_mobile_id" json:"logo_mobile"`
	FavIconID         *uint          `gorm:"column:fav_icon_id;default:NULL" json:"fav_icon_id"`
	FavIcon           *Medium        `gorm:"foreignKey:fav_icon_id" json:"fav_icon"`
	MobileIconID      *uint          `gorm:"column:mobile_icon_id;default:NULL" json:"mobile_icon_id"`
	MobileIcon        *Medium        `gorm:"foreignKey:mobile_icon_id" json:"mobile_icon"`
	VerificationCodes postgres.Jsonb `gorm:"column:verification_codes" json:"verification_codes" swaggertype:"primitive,string"`
	SocialMediaURLs   postgres.Jsonb `gorm:"column:social_media_urls" json:"social_media_urls" swaggertype:"primitive,string"`
	ContactInfo       postgres.Jsonb `gorm:"column:contact_info" json:"contact_info" swaggertype:"primitive,string"`
	Analytics         postgres.Jsonb `gorm:"column:analytics" json:"analytics" swaggertype:"primitive,string"`
	HeaderCode        string         `gorm:"column:header_code" json:"header_code"`
	FooterCode        string         `gorm:"column:footer_code" json:"footer_code"`
}

type Medium struct {
	Base
	Name        string         `gorm:"column:name" json:"name"`
	Slug        string         `gorm:"column:slug" json:"slug"`
	Type        string         `gorm:"column:type" json:"type"`
	Title       string         `gorm:"column:title" json:"title"`
	Description string         `gorm:"column:description" json:"description"`
	Caption     string         `gorm:"column:caption" json:"caption"`
	AltText     string         `gorm:"column:alt_text" json:"alt_text"`
	FileSize    int64          `gorm:"column:file_size" json:"file_size"`
	URL         postgres.Jsonb `gorm:"column:url" json:"url" swaggertype:"primitive,string"`
	Dimensions  string         `gorm:"column:dimensions" json:"dimensions"`
	MetaFields  postgres.Jsonb `gorm:"column:meta_fields" json:"meta_fields" swaggertype:"primitive,string"`
	SpaceID     uint           `gorm:"column:space_id" json:"space_id"`
}

type Claim struct {
	Base
	Claim           string         `gorm:"column:claim" json:"claim"`
	Slug            string         `gorm:"column:slug" json:"slug"`
	ClaimDate       *time.Time     `gorm:"column:claim_date" json:"claim_date" sql:"DEFAULT:NULL"`
	CheckedDate     *time.Time     `gorm:"column:checked_date" json:"checked_date" sql:"DEFAULT:NULL"`
	ClaimSources    postgres.Jsonb `gorm:"column:claim_sources" json:"claim_sources" swaggertype:"primitive,string"`
	Description     postgres.Jsonb `gorm:"column:description" json:"description" swaggertype:"primitive,string"`
	DescriptionHTML string         `gorm:"column:description_html" json:"description_html,omitempty"`
	ClaimantID      uint           `gorm:"column:claimant_id" json:"claimant_id"`
	Claimant        Claimant       `json:"claimant"`
	RatingID        uint           `gorm:"column:rating_id" json:"rating_id"`
	Rating          Rating         `json:"rating"`
	MediumID        *uint          `gorm:"column:medium_id;default:NULL" json:"medium_id"`
	Medium          *Medium        `json:"medium"`
	Fact            string         `gorm:"column:fact" json:"fact"`
	ReviewSources   postgres.Jsonb `gorm:"column:review_sources" json:"review_sources" swaggertype:"primitive,string"`
	MetaFields      postgres.Jsonb `gorm:"column:meta_fields" json:"meta_fields" swaggertype:"primitive,string"`
	SpaceID         uint           `gorm:"column:space_id" json:"space_id"`
	VideoID         *uint          `gorm:"column:video_id" json:"video_id"`
	Video           *Video         `json:"video"`
	EndTime         int            `gorm:"column:end_time" json:"end_time"`
	StartTime       int            `gorm:"column:start_time" json:"start_time"`
	Meta            postgres.Jsonb `gorm:"column:meta" json:"meta" swaggertype:"primitive,string"`
	HeaderCode      string         `gorm:"column:header_code" json:"header_code"`
	FooterCode      string         `gorm:"column:footer_code" json:"footer_code"`
}

type Claimant struct {
	config.Base
	Name            string         `gorm:"column:name" json:"name"`
	Slug            string         `gorm:"column:slug" json:"slug"`
	Description     postgres.Jsonb `gorm:"column:description" json:"description" swaggertype:"primitive,string"`
	DescriptionHTML string         `gorm:"column:description_html" json:"description_html,omitempty"`
	IsFeatured      bool           `gorm:"column:is_featured" json:"is_featured"`
	TagLine         string         `gorm:"column:tag_line" json:"tag_line"`
	MediumID        *uint          `gorm:"column:medium_id;default:NULL" json:"medium_id"`
	Medium          *model.Medium  `json:"medium"`
	MetaFields      postgres.Jsonb `gorm:"column:meta_fields" json:"meta_fields" swaggertype:"primitive,string"`
	SpaceID         uint           `gorm:"column:space_id" json:"space_id"`
	Meta            postgres.Jsonb `gorm:"column:meta" json:"meta" swaggertype:"primitive,string"`
	HeaderCode      string         `gorm:"column:header_code" json:"header_code"`
	FooterCode      string         `gorm:"column:footer_code" json:"footer_code"`
}

type Video struct {
	Base
	URL           string         `gorm:"column:url;not null" json:"url"`
	Title         string         `gorm:"column:title;not null" json:"title"`
	Slug          string         `gorm:"column:slug" json:"slug"`
	Summary       string         `gorm:"column:summary" json:"summary"`
	VideoType     string         `gorm:"column:video_type" json:"video_type"`
	SpaceID       uint           `gorm:"column:space_id" json:"space_id"`
	Status        string         `gorm:"status" json:"status"`
	TotalDuration int            `gorm:"total_duration" json:"total_duration"`
	ThumbnailURL  string         `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	PublishedDate *time.Time     `gorm:"column:published_date" json:"published_date"`
	Schemas       postgres.Jsonb `gorm:"column:schemas" json:"schemas" swaggertype:"primitive,string"`
}

type Post struct {
	Base
	Title            string         `gorm:"column:title" json:"title"`
	Subtitle         string         `gorm:"column:subtitle" json:"subtitle"`
	Slug             string         `gorm:"column:slug" json:"slug"`
	Status           string         `gorm:"column:status" json:"status"`
	IsPage           bool           `gorm:"column:is_page" json:"is_page"`
	Excerpt          string         `gorm:"column:excerpt" json:"excerpt"`
	Description      postgres.Jsonb `gorm:"column:description" json:"description" sql:"jsonb" swaggertype:"primitive,string"`
	DescriptionHTML  string         `gorm:"column:description_html" json:"description_html,omitempty"`
	IsFeatured       bool           `gorm:"column:is_featured" json:"is_featured"`
	IsSticky         bool           `gorm:"column:is_sticky" json:"is_sticky"`
	IsHighlighted    bool           `gorm:"column:is_highlighted" json:"is_highlighted"`
	FeaturedMediumID *uint          `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	Medium           *Medium        `gorm:"foreignKey:featured_medium_id" json:"medium"`
	FormatID         uint           `gorm:"column:format_id" json:"format_id" sql:"DEFAULT:NULL"`
	PublishedDate    *time.Time     `gorm:"column:published_date" json:"published_date"`
	SpaceID          uint           `gorm:"column:space_id" json:"space_id"`
	Schemas          postgres.Jsonb `gorm:"column:schemas" json:"schemas" swaggertype:"primitive,string"`
	Meta             postgres.Jsonb `gorm:"column:meta" json:"meta" swaggertype:"primitive,string"`
	HeaderCode       string         `gorm:"column:header_code" json:"header_code"`
	FooterCode       string         `gorm:"column:footer_code" json:"footer_code"`
	MetaFields       postgres.Jsonb `gorm:"column:meta_fields" json:"meta_fields" swaggertype:"primitive,string"`
}
