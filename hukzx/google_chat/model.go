package googlechat

type Header struct {
	Title      string `json:"title,omitempty"`
	Subtitle   string `json:"subtitle,omitempty"`
	ImageUrl   string `json:"imageUrl,omitempty"`
	ImageStyle string `json:"imageStyle,omitempty"`
}

type TextParagraph struct {
	Text string `json:"text,omitempty"`
}

type TextParagraphWidget struct {
	TextParagraph TextParagraph `json:"textParagraph,omitempty"`
}

type OpenLink struct {
	URL string `json:"url,omitempty"`
}
type OnClick struct {
	OpenLink OpenLink `json:"openLink,omitempty"`
}

type Button struct {
	Text    string  `json:"text,omitempty"`
	OnClick OnClick `json:"onClick,omitempty"`
}
type ButtonWidget struct {
	Button Button `json:"button,omitempty"`
}

type Image struct {
	ImageURL string  `json:"imageUrl,omitempty"`
	OnClick  OnClick `json:"onClick,omitempty"`
}

type ImageWidget struct {
	Image Image `json:"image,omitempty"`
}

type Section struct {
	Widgets []interface{} `json:"widgets,omitempty"`
}

type Card struct {
	Header   *Header   `json:"header,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

type Message struct {
	Cards []Card `json:"cards,omitempty"`
}
