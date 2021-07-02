package slack

import (
	"encoding/json"
	"fmt"
	"strings"

	coreModel "github.com/factly/dega-server/service/core/model"
	factcheckModel "github.com/factly/dega-server/service/fact-check/model"
	podcastModel "github.com/factly/dega-server/service/podcast/model"
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
		if err := json.Unmarshal(byteData, &post); err != nil {
			return nil, err
		}
		return PostToMessage(event, post)

	case "format":
		fmt := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &fmt); err != nil {
			return nil, err
		}
		return OthToMessage(entityType, event, fmt)

	case "tag":
		tag := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &tag); err != nil {
			return nil, err
		}
		return OthToMessage(entityType, event, tag)

	case "category":
		cat := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &cat); err != nil {
			return nil, err
		}
		return OthToMessage(entityType, event, cat)

	case "rating":
		rat := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &rat); err != nil {
			return nil, err
		}
		return OthToMessage(entityType, event, rat)

	case "claimant":
		claimant := map[string]interface{}{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &claimant); err != nil {
			return nil, err
		}
		return OthToMessage(entityType, event, claimant)

	case "claim":
		claim := factcheckModel.Claim{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &claim); err != nil {
			return nil, err
		}
		return ClaimToMessage(event, claim)

	case "policy":
		pol := coreModel.Policy{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &pol); err != nil {
			return nil, err
		}
		return PolicyToMessage(event, pol)

	case "podcast":
		pod := podcastModel.Podcast{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &pod); err != nil {
			return nil, err
		}
		return PodcastToMessage(event, pod)

	case "episode":
		epi := podcastModel.Episode{}
		byteData, _ := json.Marshal(whData.Payload)
		if err := json.Unmarshal(byteData, &epi); err != nil {
			return nil, err
		}
		return EpisodeToMessage(event, epi)

	default:
		return DefaultMessage(whData), nil

	}
}

func DefaultMessage(data interface{}) *Message {
	bytes, _ := json.Marshal(data)

	message := &Message{}
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "plain_text",
			Text: string(bytes),
		},
	})
	return message
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
	if post.Medium != nil {
		alt_text := "alt text"
		if post.Medium.AltText != "" {
			alt_text = post.Medium.AltText
		} else if post.Medium.Name != "" {
			alt_text = post.Medium.Name
		} else if post.Medium.Title != "" {
			alt_text = post.Medium.Title
		}
		message.Blocks = append(message.Blocks, Block{
			Type: "image",
			Title: TextBlock{
				Type: "plain_text",
				Text: fmt.Sprintf("%v", post.Title),
			},
			ImageURL: mediumURL,
			AltText:  alt_text,
		})
	}

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

func ClaimToMessage(event string, claim factcheckModel.Claim) (*Message, error) {
	message := Message{}

	// title block
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%v Claim: %v*", strings.Title(event), claim.Claim),
		},
		Accessory: ImageAccessory{
			Type:     "image",
			ImageURL: "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
			AltText:  "Factly",
		},
	})

	// fact section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: html2text.HTML2Text(claim.Fact),
		},
	})

	// date section
	dateStr := ""
	if claim.ClaimDate != nil {
		dateStr = fmt.Sprint("*Claim Date:* ", claim.ClaimDate.Format("January 2, 2006"), "\n")
	}
	if claim.CheckedDate != nil {
		dateStr = fmt.Sprint(dateStr, "*Checked Date:* ", claim.CheckedDate.Format("January 2, 2006"))
	}
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: dateStr,
		},
	})

	// claimant section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprint("*Claimant:* ", claim.Claimant.Name, "\n", "*Description:* ", html2text.HTML2Text(claim.Claimant.HTMLDescription)),
		},
	})

	// rating section
	ratingString := fmt.Sprint("*Rating:* ", claim.Rating.Name, "(", claim.Rating.NumericValue, ")", "\n", "*Description:* ", html2text.HTML2Text(claim.Rating.HTMLDescription))
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: html2text.HTML2Text(ratingString),
		},
	})

	// description section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprint("*Description:* ", html2text.HTML2Text(claim.HTMLDescription)),
		},
	})
	return &message, nil
}

func PolicyToMessage(event string, pol coreModel.Policy) (*Message, error) {
	message := Message{}

	// title block
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%v Policy: %v*", strings.Title(event), pol.Name),
		},
		Accessory: ImageAccessory{
			Type:     "image",
			ImageURL: "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
			AltText:  "Factly",
		},
	})

	// users section
	if len(pol.Users) > 0 {
		userString := "*Users:* \n"
		for _, each := range pol.Users {
			userString = fmt.Sprint(userString, each.FirstName, " ", each.LastName, " (", each.Email, ")\n")
		}
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: userString,
			},
		})
	}

	// permission section
	if len(pol.Permissions) > 0 {
		perString := "*Permissions:* \n"
		for _, each := range pol.Permissions {
			perString = fmt.Sprint(perString, "Resource: ", each.Resource)

			actionStr := "["
			for _, act := range each.Actions {
				actionStr = fmt.Sprint(actionStr, act, ",")
			}
			actionStr = strings.TrimRight(actionStr, ",")
			actionStr = fmt.Sprint(actionStr, "]")

			perString = fmt.Sprint(perString, " ", actionStr, "\n")
		}
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: perString,
			},
		})
	}
	return &message, nil
}

func PodcastToMessage(event string, podcast podcastModel.Podcast) (*Message, error) {
	message := Message{}

	mediumURL := ""
	if podcast.Medium != nil {
		urlObj := map[string]interface{}{}
		if err := json.Unmarshal(podcast.Medium.URL.RawMessage, &urlObj); err != nil {
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
			Text: fmt.Sprintf("*%v Podcast: %v* \n*Language:* %v", strings.Title(event), podcast.Title, podcast.Language),
		},
		Accessory: ImageAccessory{
			Type:     "image",
			ImageURL: "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
			AltText:  "Factly",
		},
	})

	if podcast.Medium != nil {
		// featured medium block
		alt_text := "alt text"
		if podcast.Medium.AltText != "" {
			alt_text = podcast.Medium.AltText
		} else if podcast.Medium.Name != "" {
			alt_text = podcast.Medium.Name
		} else if podcast.Medium.Title != "" {
			alt_text = podcast.Medium.Title
		}
		message.Blocks = append(message.Blocks, Block{
			Type: "image",
			Title: TextBlock{
				Type: "plain_text",
				Text: fmt.Sprintf("%v", podcast.Title),
			},
			ImageURL: mediumURL,
			AltText:  alt_text,
		})
	}

	// primary category section
	if podcast.PrimaryCategory != nil {
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: fmt.Sprint("*Primary Category:* ", podcast.PrimaryCategory.Name),
			},
		})
	}

	// categories section
	if len(podcast.Categories) > 0 {
		catString := "*Categories:* "
		for _, each := range podcast.Categories {
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

	// description section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprint("*Description:* ", html2text.HTML2Text(podcast.HTMLDescription)),
		},
	})
	return &message, nil
}

func EpisodeToMessage(event string, episode podcastModel.Episode) (*Message, error) {
	message := Message{}

	mediumURL := ""
	if episode.Medium != nil {
		urlObj := map[string]interface{}{}
		if err := json.Unmarshal(episode.Medium.URL.RawMessage, &urlObj); err != nil {
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
			Text: fmt.Sprintf("*%v Episode: %v* \n*Season:* %v  *Episode:* %v", strings.Title(event), episode.Title, episode.Season, episode.Episode),
		},
		Accessory: ImageAccessory{
			Type:     "image",
			ImageURL: "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
			AltText:  "Factly",
		},
	})

	// featured medium block
	if episode.Medium != nil {
		alt_text := "alt text"
		if episode.Medium.AltText != "" {
			alt_text = episode.Medium.AltText
		} else if episode.Medium.Name != "" {
			alt_text = episode.Medium.Name
		} else if episode.Medium.Title != "" {
			alt_text = episode.Medium.Title
		}
		message.Blocks = append(message.Blocks, Block{
			Type: "image",
			Title: TextBlock{
				Type: "plain_text",
				Text: fmt.Sprintf("%v", episode.Title),
			},
			ImageURL: mediumURL,
			AltText:  alt_text,
		})
	}

	// podcast section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprint("*Podcast:* ", episode.Podcast.Title),
		},
	})

	// published date section
	if episode.PublishedDate != nil {
		message.Blocks = append(message.Blocks, Block{
			Type: "section",
			Text: TextBlock{
				Type: "mrkdwn",
				Text: fmt.Sprint("*Published Date:* ", episode.PublishedDate.Format("January 2, 2006")),
			},
		})
	}

	// description section
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprint("*Description* ", html2text.HTML2Text(episode.HTMLDescription)),
		},
	})

	// audio section
	audiourl := fmt.Sprintf(`<a href="%v">%v</a>`, episode.AudioURL, episode.Title)
	message.Blocks = append(message.Blocks, Block{
		Type: "section",
		Text: TextBlock{
			Type: "mrkdwn",
			Text: fmt.Sprint(`*Audio URL:* `, html2text.HTML2Text(audiourl)),
		},
	})
	return &message, nil
}
