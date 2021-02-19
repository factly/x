package editorx

import (
	"fmt"
)

func getHeaderElement(text string, level int) string {
	return fmt.Sprintf("\n<h%v>%s</h%v>", level, text, level)
}

func getListElement(items []interface{}, tpe string) string {
	itemsHTML := ""
	for _, item := range items {
		itemStr := item.(string)
		itemsHTML = fmt.Sprintf("%s\n<li>%s</li>", itemsHTML, itemStr)
	}

	if tpe == "ordered" {
		return fmt.Sprintf("\n<ol>%s\n</ol>", itemsHTML)
	} else if tpe == "unordered" {
		return fmt.Sprintf("\n<ul>%s\n</ul>", itemsHTML)
	}
	return ""
}

func getTableElement(content []interface{}) string {
	html := "\n<table>"

	for i, r := range content {
		row := r.([]interface{})
		html = fmt.Sprint(html, "\n<tr>")
		for _, c := range row {
			cstr := c.(string)
			if i == 0 {
				html = fmt.Sprintf("%s\n<th>%s</th>", html, cstr)
			} else {
				html = fmt.Sprintf("%s\n<td>%s</td>", html, cstr)
			}
		}
		html = fmt.Sprint(html, "\n</tr>")
	}
	html = fmt.Sprint(html, "\n</table>")

	return html
}

func getCodeElement(code string) string {
	html := "\n<pre><code style=\"display:block;\">"
	html = fmt.Sprint(html, "\n", code, "\n</code></pre>")
	return html
}

func getImageElement(image map[string]interface{}) string {
	urlMap := image["url"].(map[string]interface{})

	var url string
	if urlint, found := urlMap["proxy"]; found {
		url = urlint.(string)
	} else {
		url = urlMap["raw"].(string)
	}

	altText := image["alt_text"].(string)
	caption := image["caption"].(string)
	slug := image["slug"].(string)

	html := fmt.Sprintf("\n<div class=\"image\"><img src=\"%s\" id=\"%s\" alt=\"%s\"> </div>", url, slug, altText)
	return fmt.Sprint(html, "\n<p>", caption, "</p>")
}
