package slack

import (
	"encoding/json"
	"fmt"
	"strings"

	whmodel "github.com/factly/hukz/model"
	"github.com/factly/x/hukzx"
	"github.com/k3a/html2text"
)

func ToMessage(whData whmodel.WebhookData) (*Message, error) {
	eventTokens := strings.Split(whData.Event, ".")
	entityType := eventTokens[0]
	event := eventTokens[1]

	switch entityType {
	case "post":
		post := hukzx.Post{}
		byteData, _ := json.Marshal(whData.Payload)
		_ = json.Unmarshal(byteData, &post)
		return PostToMessage(event, post)

	case "format":
		fmt := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		_ = json.Unmarshal(byteData, &fmt)
		return OthToMessage(entityType, event, fmt)

	case "tag":
		tag := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		_ = json.Unmarshal(byteData, &tag)
		return OthToMessage(entityType, event, tag)

	case "category":
		cat := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		_ = json.Unmarshal(byteData, &cat)
		return OthToMessage(entityType, event, cat)

	case "rating":
		rat := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		_ = json.Unmarshal(byteData, &rat)
		return OthToMessage(entityType, event, rat)

	case "claimant":
		claimant := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		_ = json.Unmarshal(byteData, &claimant)
		return OthToMessage(entityType, event, claimant)

	}

	return nil, nil
}

func PostToMessage(event string, post hukzx.Post) (*Message, error) {
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

	// title block
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%v Post: %v* \n%v", strings.Title(event), post.Title, post.Excerpt),
		},
		Accessory: ImageAccessory{
			Type:     "image",
			ImageURL: "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
			AltText:  "Factly",
		},
	})

	// featured medium block
	message.Blocks = append(message.Blocks, Block{
		Type: "image",
		Title: TextBlock{
			Type: "plain_text",
			Text: fmt.Sprintf("%v", post.Title),
		},
		ImageURL: mediumURL,
		AltText:  post.Medium.Title,
	})

	// published date section
	if post.PublishedDate != nil {
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*Published Date: * %v", post.PublishedDate.Format("January 2, 2006")),
			},
		})
	}

	// claims section
	if len(post.Claims) > 0 {
		claimString := ""
		for _, each := range post.Claims {
			claimString = fmt.Sprint(claimString, "\n*Claim:* ")
			claimString = fmt.Sprint(claimString, each.Claim, "\n")
			claimString = fmt.Sprint(claimString, "*Fact:* ", each.Fact, "\n")
			claimString = fmt.Sprint(claimString, "*Claimant:* ", each.Claimant.Name, "\n")
			claimString = fmt.Sprint(claimString, "*Rating:* ", each.Rating.Name, "\n")
			claimString = fmt.Sprint(claimString, "*ClaimDate:* ", each.ClaimDate.Format("01/02/2006"), "\n")
			claimString = fmt.Sprint(claimString, "*Description:* ", html2text.HTML2Text(each.HTMLDescription), "\n")
		}
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: claimString,
			},
		})
	}

	// categories section
	if len(post.Categories) > 0 {
		catString := "*Categories:* "
		for _, each := range post.Categories {
			catString = fmt.Sprint(catString, each.Name, ", ")
		}
		catString = strings.TrimRight(catString, ", ")
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: catString,
			},
		})
	}

	// tags section
	if len(post.Categories) > 0 {
		tagString := "*Tags:* "
		for _, each := range post.Tags {
			tagString = fmt.Sprint(tagString, each.Name, ", ")
		}
		tagString = strings.TrimRight(tagString, ", ")
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: tagString,
			},
		})
	}

	// description section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: html2text.HTML2Text(post.HTMLDescription),
		},
	})

	return &message, nil
}

func OthToMessage(entityType, action string, obj map[string]interface{}) (*Message, error) {
	message := Message{}
	name := obj["name"].(string)
	var desc string
	if entityType != "format" {
		if in, ok := obj["html_description"]; ok && in != nil {
			desc = obj["html_description"].(string)
		}
	} else {
		if in, ok := obj["description"]; ok && in != nil {
			desc = obj["description"].(string)
		}
	}

	// title block
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%v %v: %v*", strings.Title(action), strings.Title(entityType), name),
		},
		Accessory: ImageAccessory{
			Type:     "image",
			ImageURL: "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
			AltText:  "Factly",
		},
	})

	// description section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: html2text.HTML2Text(desc),
		},
	})

	return &message, nil
}
