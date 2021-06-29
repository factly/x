package hukzx

import (
	coreModel "github.com/factly/dega-server/service/core/model"
	factCheckModel "github.com/factly/dega-server/service/fact-check/model"
)

type Post struct {
	coreModel.Post
	Authors []coreModel.Author     `json:"authors"`
	Claims  []factCheckModel.Claim `json:"claims"`
}
