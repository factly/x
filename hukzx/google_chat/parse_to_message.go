package googlechat

import (
	"encoding/json"
	"fmt"
	"strings"

	whmodel "github.com/factly/hukz/model"
	"github.com/factly/x/hukzx"
)

func ToMessage(whData whmodel.WebhookData) (*Message, error) {
	entityType := strings.Split(whData.Event, ".")[0]

	switch entityType {
	case "post":
		post := hukzx.Post{}
		byteData, _ := json.Marshal(whData.Payload)

		_ = json.Unmarshal(byteData, &post)

		return PostToMessage(post)

	}

	return nil, nil
}

func PostToMessage(post hukzx.Post) (*Message, error) {
	message := Message{}

	mediumURL := ""
	if post.Medium != nil {
		urlObj := map[string]interface{}{}
		if err := json.Unmarshal(post.Medium.URL.RawMessage, &urlObj); err != nil {
			return nil, err
		}
		if _, ok := urlObj["raw"]; ok {
			mediumURL = urlObj["raw"].(string)
		}
	}

	postCard := Card{}

	if post.Medium != nil {
		postCard.Header = &Header{
			Title:      post.Title,
			Subtitle:   post.Excerpt,
			ImageUrl:   mediumURL,
			ImageStyle: "IMAGE",
		}
	} else {
		postCard.Header = &Header{
			Title:    post.Title,
			Subtitle: post.Excerpt,
		}
	}

	// featured medium section
	fmSection := Section{}
	imgWidget := ImageWidget{
		Image: Image{
			ImageURL: mediumURL,
		},
	}
	fmSection.Widgets = append(fmSection.Widgets, imgWidget)
	postCard.Sections = append(postCard.Sections, fmSection)

	// published date section
	if post.PublishedDate != nil {
		dateSection := Section{}
		dateTxtWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: fmt.Sprint("<b>Published Date:</b> ", post.PublishedDate.Format("January 2, 2006")),
			},
		}
		dateSection.Widgets = append(dateSection.Widgets, dateTxtWidget)
		postCard.Sections = append(postCard.Sections, dateSection)
	}

	// claims section
	if len(post.Claims) > 0 {
		claimSection := Section{}
		claimString := ""
		for _, each := range post.Claims {
			claimString = "<b>Claim: </b>"
			claimString = fmt.Sprint(claimString, each.Claim, "<br>")
			claimString = fmt.Sprint(claimString, "<b>Fact: </b>", each.Fact, "<br>")
			claimString = fmt.Sprint(claimString, "<b>Claimant: </b>", each.Claimant.Name, "<br>")
			claimString = fmt.Sprint(claimString, "<b>Rating: </b>", each.Rating.Name, "<br>")
			claimString = fmt.Sprint(claimString, "<b>ClaimDate: </b>", each.ClaimDate.Format("01/02/2006"), "<br>")
			claimString = fmt.Sprint(claimString, "<b>Description: </b>", each.HTMLDescription, "<br>")
			claimWidget := TextParagraphWidget{
				TextParagraph: TextParagraph{
					Text: claimString,
				},
			}
			claimSection.Widgets = append(claimSection.Widgets, claimWidget)
		}
		postCard.Sections = append(postCard.Sections, claimSection)
	}

	// categories section
	if len(post.Categories) > 0 {
		catSection := Section{}
		catString := "<b>Categories: </b>"
		for _, each := range post.Categories {
			catString = fmt.Sprint(catString, each.Name, ", ")
		}
		catString = strings.TrimRight(catString, ", ")
		catWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: catString,
			},
		}
		catSection.Widgets = append(catSection.Widgets, catWidget)
		postCard.Sections = append(postCard.Sections, catSection)
	}

	// tags section
	if len(post.Tags) > 0 {
		tagSection := Section{}
		tagString := "<b>Tags: </b>"
		for _, each := range post.Tags {
			tagString = fmt.Sprint(tagString, each.Name, ", ")
		}
		tagString = strings.TrimRight(tagString, ", ")
		tagWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: tagString,
			},
		}
		tagSection.Widgets = append(tagSection.Widgets, tagWidget)
		postCard.Sections = append(postCard.Sections, tagSection)
	}

	// description section
	descSection := Section{}
	descriptionWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: post.HTMLDescription,
		},
	}
	descSection.Widgets = append(descSection.Widgets, descriptionWidget)
	postCard.Sections = append(postCard.Sections, descSection)

	message.Cards = append(message.Cards, postCard)
	return &message, nil
}
