package editorx

import (
	"errors"
	"fmt"
)

// EditorjsToHTML converts editorjs description blocks into html
func EditorjsToHTML(raw map[string]interface{}) (string, error) {
	html := ""

	blocks, ok := raw["blocks"].([]interface{})
	if !ok {
		return "", errors.New("type error for blocks")
	}

	for i, blk := range blocks {
		block := blk.(map[string]interface{})

		btype, ok := block["type"].(string)
		if !ok {
			return "", errors.New(fmt.Sprint("type error for type in block #", i))
		}
		bdata, ok := block["data"].(map[string]interface{})
		if !ok {
			return "", errors.New(fmt.Sprint("type error for data in block #", i))
		}

		if btype == "header" {
			headerText := bdata["text"].(string)
			headerLevel := bdata["level"].(float64)
			header := getHeaderElement(headerText, int(headerLevel))
			html = fmt.Sprint(html, header)

		} else if btype == "paragraph" {
			paraText := bdata["text"].(string)
			html = fmt.Sprintf("%s\n<p>%s</p>", html, paraText)

		} else if btype == "list" {
			listStyle := bdata["style"].(string)
			listItems := bdata["items"].([]interface{})
			html = fmt.Sprint(html, getListElement(listItems, listStyle))

		} else if btype == "quote" {
			text := bdata["text"].(string)
			html = fmt.Sprintf("%s\n<blockquote>%s</blockquote>", html, text)

		} else if btype == "raw" {
			text := bdata["html"].(string)
			html = fmt.Sprintf("%s\n%s", html, text)

		} else if btype == "table" {
			content := bdata["content"].([]interface{})
			html = fmt.Sprint(html, getTableElement(content))

		} else if btype == "code" {
			code := bdata["code"].(string)
			html = fmt.Sprint(html, getCodeElement(code))

		} else if btype == "delimiter" {
			html = fmt.Sprint(html, "\n<hr>")

		} else if btype == "uppy" {
			imageNodes := bdata["nodes"].([]interface{})
			for _, each := range imageNodes {
				image := each.(map[string]interface{})
				html = fmt.Sprint(html, getImageElement(image))
			}

		} else if btype == "embed" {
			embedHTML, ok := bdata["html"].(string)
			if ok {
				html = fmt.Sprint(html, "\n", embedHTML)
			}
		}
	}
	return html, nil
}
