package editorx

import (
	"bytes"
	"html/template"
	"strings"
	"time"
)

var BasePath = "templates"

// EditorjsToHTML converts editorjs description blocks into html
func EditorjsToHTML(raw map[string]interface{}) (string, error) {

	tpl := SetupTemplates(BasePath)

	bmap, err := BlockMap(raw)
	if err != nil {
		return "", err
	}

	if err = CheckBlocks(bmap); err != nil {
		return "", err
	}

	var tplBuff bytes.Buffer
	err = tpl.ExecuteTemplate(&tplBuff, "description.gohtml", bmap)
	if err != nil {
		return "", err
	}

	html := strings.TrimSpace(tplBuff.String())
	html = strings.ReplaceAll(html, "\n", "")
	html = strings.ReplaceAll(html, "\t", "")

	return html, nil
}

// SetupTemplates setups the templates
func SetupTemplates(basePath string) *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"bmap":         BlockMap,
		"dateFmt":      formatDate,
		"dateVal":      validateDate,
		"noesc":        noescape,
		"multipleImgs": multipleUppy,
	}).ParseGlob(basePath + "/*"))
}

func formatDate(date time.Time) string {
	return date.Format("01/02/2006")
}

func validateDate(date time.Time) bool {
	return !date.Equal(time.Time{})
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func multipleUppy(uppyData map[string]interface{}) bool {
	if _, found := uppyData["nodes"]; !found {
		return false
	}
	return true
}
