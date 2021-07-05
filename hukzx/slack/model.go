package slack

type Block struct {
	Type      string      `json:"type,omitempty"`
	Text      interface{} `json:"text,omitempty"`
	Title     interface{} `json:"title,omitempty"`
	Accessory interface{} `json:"accessory,omitempty"`
	ImageURL  string      `json:"image_url,omitempty"`
	AltText   string      `json:"alt_text,omitempty"`
}

type TextBlock struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

type ImageAccessory struct {
	Type     string `json:"type,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	AltText  string `json:"alt_text,omitempty"`
}

type Message struct {
	Blocks []Block `json:"blocks,omitempty"`
}
