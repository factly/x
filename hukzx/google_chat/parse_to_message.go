package googlechat

import (
	"encoding/json"
	"fmt"
	"strings"

	coreModel "github.com/factly/dega-server/service/core/model"
	factcheckModel "github.com/factly/dega-server/service/fact-check/model"
	podcastModel "github.com/factly/dega-server/service/podcast/model"
	whmodel "github.com/factly/hukz/model"
	"github.com/factly/x/hukzx"
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
	message := Message{}
	card := Card{}
	sec := Section{}
	txtWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: string(bytes),
		},
	}
	sec.Widgets = append(sec.Widgets, txtWidget)
	card.Sections = append(card.Sections, sec)
	message.Cards = append(message.Cards, card)

	return &message
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

	postCard := Card{}

	postCard.Header = &Header{
		Title:      fmt.Sprint(strings.Title(event), " Post: ", post.Title),
		Subtitle:   post.Excerpt,
		ImageUrl:   "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		ImageStyle: "IMAGE",
	}

	// featured medium section
	if post.Medium != nil {
		fmSection := Section{}
		imgWidget := ImageWidget{
			Image: Image{
				ImageURL: mediumURL,
			},
		}
		fmSection.Widgets = append(fmSection.Widgets, imgWidget)
		postCard.Sections = append(postCard.Sections, fmSection)
	}

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
	card := Card{}

	card.Header = &Header{
		Title:      fmt.Sprint(strings.Title(entityType), " '", name, "' ", action),
		ImageUrl:   "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		ImageStyle: "IMAGE",
	}

	// description section
	descSection := Section{}
	descriptionWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: desc,
		},
	}
	descSection.Widgets = append(descSection.Widgets, descriptionWidget)
	card.Sections = append(card.Sections, descSection)

	message.Cards = append(message.Cards, card)
	return &message, nil
}

func ClaimToMessage(event string, claim factcheckModel.Claim) (*Message, error) {
	message := Message{}

	card := Card{}

	card.Header = &Header{
		Title:      fmt.Sprint(strings.Title(event), " Claim: ", claim.Claim),
		ImageUrl:   "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		ImageStyle: "IMAGE",
	}

	// fact section
	factSec := Section{}
	factWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: fmt.Sprint("<b>Fact: </b>", claim.Fact),
		},
	}
	factSec.Widgets = append(factSec.Widgets, factWidget)
	card.Sections = append(card.Sections, factSec)

	// date section
	dateSec := Section{}
	dateStr := ""
	if claim.ClaimDate != nil {
		dateStr = fmt.Sprint("<b>Claim Date: </b>", claim.ClaimDate.Format("January 2, 2006"), "<br>")
	}
	if claim.CheckedDate != nil {
		dateStr = fmt.Sprint(dateStr, "<b>Checked Date: </b>", claim.CheckedDate.Format("January 2, 2006"))
	}
	dateWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: dateStr,
		},
	}
	dateSec.Widgets = append(dateSec.Widgets, dateWidget)
	card.Sections = append(card.Sections, dateSec)

	// claimant section
	claimantSec := Section{}
	claimantStr := fmt.Sprint("<b>Claimant: </b>", claim.Claimant.Name, "<br>", "<b>Description: </b>", claim.Claimant.HTMLDescription)
	claimantWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: claimantStr,
		},
	}
	claimantSec.Widgets = append(claimantSec.Widgets, claimantWidget)
	card.Sections = append(card.Sections, claimantSec)

	// rating section
	ratingSec := Section{}
	ratingCol := claim.Rating.BackgroundColour
	colorObj := map[string]interface{}{}

	_ = json.Unmarshal(ratingCol.RawMessage, &colorObj)
	color := ""
	if col, ok := colorObj["hex"]; ok && col != nil {
		color = colorObj["hex"].(string)
	}

	ratingString := fmt.Sprintf("<b>Rating: </b> <font color=\"%v\">%v (%v)</font> <br> <b>Description: </b>%v", color, claim.Rating.Name, claim.Rating.NumericValue, claim.Rating.HTMLDescription)

	ratingWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: ratingString,
		},
	}
	ratingSec.Widgets = append(ratingSec.Widgets, ratingWidget)
	card.Sections = append(card.Sections, ratingSec)

	// description section
	descSection := Section{}
	descriptionWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: fmt.Sprint("<b>Description: </b>", claim.HTMLDescription),
		},
	}
	descSection.Widgets = append(descSection.Widgets, descriptionWidget)
	card.Sections = append(card.Sections, descSection)

	message.Cards = append(message.Cards, card)
	return &message, nil
}

func PolicyToMessage(event string, pol coreModel.Policy) (*Message, error) {
	message := Message{}

	card := Card{}

	card.Header = &Header{
		Title:      fmt.Sprint(strings.Title(event), " Policy: ", pol.Name),
		Subtitle:   pol.Description,
		ImageUrl:   "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		ImageStyle: "IMAGE",
	}

	// users section
	if len(pol.Users) > 0 {
		userSection := Section{}
		userString := "<b>Users: </b> <br>"
		for _, each := range pol.Users {
			userString = fmt.Sprint(userString, each.FirstName, " ", each.LastName, " (", each.Email, ")<br>")
		}
		userWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: userString,
			},
		}
		userSection.Widgets = append(userSection.Widgets, userWidget)
		card.Sections = append(card.Sections, userSection)
	}

	// permission section
	if len(pol.Permissions) > 0 {
		perSection := Section{}
		perString := "<b>Permissions: </b> <br>"
		for _, each := range pol.Permissions {
			perString = fmt.Sprint(perString, "Resource: ", each.Resource)

			actionStr := "["
			for _, act := range each.Actions {
				actionStr = fmt.Sprint(actionStr, act, ",")
			}
			actionStr = strings.TrimRight(actionStr, ",")
			actionStr = fmt.Sprint(actionStr, "]")

			perString = fmt.Sprint(perString, " ", actionStr, "<br>")
		}
		perWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: perString,
			},
		}
		perSection.Widgets = append(perSection.Widgets, perWidget)
		card.Sections = append(card.Sections, perSection)
	}
	message.Cards = append(message.Cards, card)
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

	card := Card{}

	card.Header = &Header{
		Title:      fmt.Sprint(strings.Title(event), " Podcast: ", podcast.Title),
		Subtitle:   fmt.Sprint("Language: ", podcast.Language),
		ImageUrl:   "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		ImageStyle: "IMAGE",
	}

	// featured medium section
	if podcast.Medium != nil {
		fmSection := Section{}
		imgWidget := ImageWidget{
			Image: Image{
				ImageURL: mediumURL,
			},
		}
		fmSection.Widgets = append(fmSection.Widgets, imgWidget)
		card.Sections = append(card.Sections, fmSection)
	}

	// primary category section
	if podcast.PrimaryCategory != nil {
		pcSection := Section{}
		pcWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: fmt.Sprint("<b>Primary Category: </b>", podcast.PrimaryCategory.Name),
			},
		}
		pcSection.Widgets = append(pcSection.Widgets, pcWidget)
		card.Sections = append(card.Sections, pcSection)
	}

	// categories section
	if len(podcast.Categories) > 0 {
		catSection := Section{}
		catString := "<b>Categories: </b>"
		for _, each := range podcast.Categories {
			catString = fmt.Sprint(catString, each.Name, ", ")
		}
		catString = strings.TrimRight(catString, ", ")
		catWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: catString,
			},
		}
		catSection.Widgets = append(catSection.Widgets, catWidget)
		card.Sections = append(card.Sections, catSection)
	}

	// description section
	descSection := Section{}
	descriptionWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: fmt.Sprint("<b>Description: </b>", podcast.HTMLDescription),
		},
	}
	descSection.Widgets = append(descSection.Widgets, descriptionWidget)
	card.Sections = append(card.Sections, descSection)

	message.Cards = append(message.Cards, card)
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

	card := Card{}

	card.Header = &Header{
		Title:      fmt.Sprint(strings.Title(event), " Episode: ", episode.Title),
		Subtitle:   fmt.Sprint("Season: ", episode.Season, " Episode: ", episode.Episode),
		ImageUrl:   "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		ImageStyle: "IMAGE",
	}

	// featured medium section
	if episode.Medium != nil {
		fmSection := Section{}
		imgWidget := ImageWidget{
			Image: Image{
				ImageURL: mediumURL,
			},
		}
		fmSection.Widgets = append(fmSection.Widgets, imgWidget)
		card.Sections = append(card.Sections, fmSection)
	}

	// podcast section
	podSection := Section{}
	podWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: fmt.Sprint("<b>Podcast: </b>", episode.Podcast.Title),
		},
	}
	podSection.Widgets = append(podSection.Widgets, podWidget)
	card.Sections = append(card.Sections, podSection)

	// published date section
	if episode.PublishedDate != nil {
		dateSection := Section{}
		dateTxtWidget := TextParagraphWidget{
			TextParagraph: TextParagraph{
				Text: fmt.Sprint("<b>Published Date:</b> ", episode.PublishedDate.Format("January 2, 2006")),
			},
		}
		dateSection.Widgets = append(dateSection.Widgets, dateTxtWidget)
		card.Sections = append(card.Sections, dateSection)
	}

	// description section
	descSection := Section{}
	descriptionWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: fmt.Sprint("<b>Description: </b>", episode.HTMLDescription),
		},
	}
	descSection.Widgets = append(descSection.Widgets, descriptionWidget)
	card.Sections = append(card.Sections, descSection)

	// audio section
	audioSection := Section{}
	audioWidget := TextParagraphWidget{
		TextParagraph: TextParagraph{
			Text: fmt.Sprintf(`<b>Audio URL: </b><a href="%v">%v</a> `, episode.AudioURL, episode.Title),
		},
	}
	audioSection.Widgets = append(audioSection.Widgets, audioWidget)
	card.Sections = append(card.Sections, audioSection)

	message.Cards = append(message.Cards, card)
	return &message, nil
}
