package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	coreModel "github.com/factly/dega-server/service/core/model"
	factcheckModel "github.com/factly/dega-server/service/fact-check/model"
	whmodel "github.com/factly/hukz/model"
	"github.com/factly/x/hukzx"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func TestParseToSlackMessage(t *testing.T) {
	t.Run("run ToMessage function for post", func(t *testing.T) {
		now := time.Now()
		mediumURL := map[string]interface{}{
			"raw": "https://factly.in/wp-content/uploads//2021/01/factly-logo-200-11.png",
		}
		urlBytes, _ := json.Marshal(mediumURL)
		message, err := ToMessage(whmodel.WebhookData{
			Event:    "post.created",
			Contains: []string{"post"},
			Payload: hukzx.Post{
				Post: coreModel.Post{
					Title:           "Test Title",
					Subtitle:        "Test Subtitle",
					Slug:            "test-title",
					PublishedDate:   &now,
					HTMLDescription: `<p></p>     <p></p>     <p>Within a year of the emergence of COVID-19, vaccines have been developed, approved, and administered to people in various parts of the world. Resources were diverted towards the development of vaccines and countries were in a race to develop and procure vaccines to curb the spread of the pandemic. Usually, it takes years of research and testing before the vaccines are tried on human beings. But the quick spread of the COVID-19 pandemic across the world left no choice but to accelerate the development, trial, and production of vaccines, in record time. As of 15 May 2021, 1.45 billion <a href="https://ourworldindata.org/covid-vaccinations?country=OWID_WRL" rel="noopener noreferrer" target="_blank">doses of vaccines</a> have been administered across the world. According to <a href="https://www.who.int/publications/m/item/draft-landscape-of-covid-19-candidate-vaccines" rel="noopener noreferrer" target="_blank">WHO’s draft landscape of novel coronavirus candidate vaccine</a> development worldwide, 100 vaccines were in clinical development and 184 vaccines were in pre-clinical development stages on 14 May 2021.</p>     <p></p>     <p></p>     <p><strong>14 Vaccines have been approved/authorized globally</strong></p>     <p></p>     <p>Globally, fourteen vaccines have been granted emergency use authorizations or have been approved by national regulatory authorities in various countries, according to the <a href="https://www.nytimes.com/interactive/2020/science/coronavirus-vaccine-tracker.html" rel="noopener noreferrer" target="_blank">New York Times Vaccine Tracker</a>. These 14 vaccines are listed below. Besides, there are 27 vaccine candidates undergoing phase-3 trials, 37 in phase-2, and 49 vaccine candidates in phase-1 trials. </p>     <p></p>        <ul>                <li>&lt;strong&gt;Vaxzevria by Oxford-AstraZeneca, UK (also known as Covishield in India)&lt;/strong&gt;</li>            </ul>         <p></p>`,
					Status:          "publish",
					Excerpt:         "This is a test post",
					Medium: &coreModel.Medium{
						Title: "Test Medium",
						URL:   postgres.Jsonb{RawMessage: urlBytes},
					},
					Categories: []coreModel.Category{
						{
							Name: "Test Post",
						},
						{
							Name: "New Post",
						},
					},
					Tags: []coreModel.Tag{
						{
							Name: "Tag1",
						},
					},
				},
				Authors: []coreModel.Author{
					{
						FirstName: "Test",
						LastName:  "User",
						Email:     "testuser@org.com",
					},
				},
				Claims: []factcheckModel.Claim{
					{
						Claim:           "This is a test claim.",
						Fact:            "This is real fact.",
						ClaimDate:       &now,
						HTMLDescription: "<h2>Test Claim</h2>",
						Claimant: factcheckModel.Claimant{
							Name: "Test Claimant",
							Slug: "test-claimant",
						},
						Rating: factcheckModel.Rating{
							Name:            "True",
							NumericValue:    1,
							HTMLDescription: "<p>True Rating</p>",
						},
					},
					{
						Claim:           "This is another test claim.",
						Fact:            "This is false fact.",
						ClaimDate:       &now,
						HTMLDescription: "<h2>A False Claim</h2>",
						Claimant: factcheckModel.Claimant{
							Name: "Test Claimant",
							Slug: "test-claimant",
						},
						Rating: factcheckModel.Rating{
							Name:            "False",
							NumericValue:    5,
							HTMLDescription: "<p>False Rating</p>",
						},
					},
				},
			},
		})

		if err != nil {
			log.Println(err.Error())
			t.Fail()
		}
		fmt.Println("MESSAGE: ", message)

	})

	t.Run("run ToMessage function for category", func(t *testing.T) {
		now := time.Now()
		message, err := ToMessage(whmodel.WebhookData{
			Event:     "category.created",
			CreatedAt: now,
			Contains:  []string{"category"},
			Payload: coreModel.Category{
				Name:            "Test Category",
				HTMLDescription: "This is test description for category",
			},
		})

		if err != nil {
			log.Println(err.Error())
			t.Fail()
		}
		fmt.Println("MESSAGE: ", message)

	})

	t.Run("run ToMessage function for format", func(t *testing.T) {
		now := time.Now()
		message, err := ToMessage(whmodel.WebhookData{
			Event:     "format.created",
			CreatedAt: now,
			Contains:  []string{"format"},
			Payload: coreModel.Format{
				Name:        "Test Format",
				Description: "This is test description for a format",
			},
		})

		if err != nil {
			log.Println(err.Error())
			t.Fail()
		}
		fmt.Println("MESSAGE: ", message)

	})
}
