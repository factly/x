package editorx

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func TestErrorx(t *testing.T) {
	descstr := `
	{
		"time": 1605767087876,
		"blocks": [
			{
				"type": "header",
				"data": {
					"text": "Test heading",
					"level": 2
				}
			},
			{
				"type": "paragraph",
				"data": {
					"text": "Test paragraph text is here"
				}
			},
			{
				"type": "list",
				"data": {
					"style": "ordered",
					"items": [
						"Ordered list item 1",
						"Ordered list item 2"
					]
				}
			},
			{
				"type": "list",
				"data": {
					"style": "unordered",
					"items": [
						"Unordered list item 1",
						"Unordered list item 2"
					]
				}
			},
			{
				"type": "quote",
				"data": {
					"text": "This is a quote from something or someone",
					"caption": "This is quote caption",
					"alignment": "left"
				}
			},
			{
				"type": "raw",
				"data": {
					"html": "<p> This is some raw html shit </p>"
				}
			},
			{
				"type": "table",
				"data": {
					"content": [
						[
							"Name",
							"Number"
						],
						[
							"Test 1",
							"1"
						],
						[
							"Test 2",
							"2"
						],
						[
							"Test 3",
							"3"
						]
					]
				}
			},
			{
				"type": "code",
				"data": {
					"code": "package main\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello world\")\n}"
				}
			},
			{
				"type": "delimiter",
				"data": {}
			},
			{
				"type": "embed",
				"data": {
					"html": "<blockquote class=\"twitter-tweet\" data-lang=\"en_US\"><p lang=\"en\" dir=\"ltr\">&quot;Sometimes I like to hold things like this and pretend I&#39;m a giant.&quot; Line by Craig Bierko. Funniest man I have ever met. <a href=\"https://t.co/b1TmUuckjH\">pic.twitter.com/b1TmUuckjH</a></p>&mdash; matthew perry (@MatthewPerry) <a href=\"https://twitter.com/MatthewPerry/status/1329015312004632576?ref_src=twsrc%5Etfw\">November 18, 2020</a></blockquote>\n<script async src=\"https://platform.twitter.com/widgets.js\" charset=\"utf-8\"></script>\n",
					"meta": {
						"author": "matthew perry",
						"author_url": "https://twitter.com/MatthewPerry",
						"canonical": "https://twitter.com/MatthewPerry/status/1329015312004632576",
						"description": "\"Sometimes I like to hold things like this and pretend I'm a giant.\" Line by Craig Bierko. Funniest man I have ever met. pic.twitter.com/b1TmUuckjH— matthew perry (@MatthewPerry) November 18, 2020\n\n",
						"site": "Twitter",
						"title": "matthew perry on Twitter"
					},
					"caption": ""
				}
			},
			{
				"type": "uppy",
				"data": {
					"total": 1,
					"nodes": [
						{
							"id": 8,
							"created_at": "2020-11-19T06:24:47.273026071Z",
							"updated_at": "2020-11-19T06:24:47.273026071Z",
							"deleted_at": {},
							"name": "test.png",
							"slug": "test-png",
							"type": "image/png",
							"title": " ",
							"description": "",
							"caption": "A Caption",
							"alt_text": "Test alt text",
							"file_size": 37257,
							"url": {
								"proxy": "http://127.0.0.1:7001/dega/test-space/2020/10/1605767086916_test.png",
								"raw": "http://localhost:9000/dega/test-space/2020/10/1605767086916_test.png"
							},
							"dimensions": "100x100",
							"space_id": 2
						}
					]
				}
			},
			{
				"type": "uppy",
				"data": {
							"id": 8,
							"created_at": "2020-11-19T06:24:47.273026071Z",
							"updated_at": "2020-11-19T06:24:47.273026071Z",
							"deleted_at": {},
							"name": "test.png",
							"slug": "test-png",
							"type": "image/png",
							"title": " ",
							"description": "",
							"caption": "A Caption",
							"alt_text": "Test alt text",
							"file_size": 37257,
							"url": {
								"proxy": "http://127.0.0.1:7001/dega/test-space/2020/10/1605767086916_test.png",
								"raw": "http://localhost:9000/dega/test-space/2020/10/1605767086916_test.png"
							},
							"dimensions": "100x100",
							"space_id": 2
						}
			}
		],
		"version": "2.19.0"
	}
	`
	testDescription := postgres.Jsonb{
		RawMessage: []byte(descstr),
	}

	outHtml := `
	<h2>Test heading</h2>     <p>Test paragraph text is here</p>        <ol>                <li>Ordered list item 1</li>                <li>Ordered list item 2</li>            </ol>            <ul>                <li>Unordered list item 1</li>                <li>Unordered list item 2</li>            </ul>        <blockquote>This is a quote from something or someone</blockquote>    <p> This is some raw html shit </p>    <table style="border: 1px solid black; width: 50%;">            <tr>                                <th> Name </td>                        <th> Number </td>                        </tr>            <tr>                                 <td>  Test 1  </th>                          <td>  1  </th>                         </tr>            <tr>                                 <td>  Test 2  </th>                          <td>  2  </th>                         </tr>            <tr>                                 <td>  Test 3  </th>                          <td>  3  </th>                         </tr>        </table>      <pre>    <code style="display:block">        package mainimport &#34;fmt&#34;func main() {    fmt.Println(&#34;Hello world&#34;)}    </code>    </pre>    <hr>    <blockquote class="twitter-tweet" data-lang="en_US"><p lang="en" dir="ltr">&quot;Sometimes I like to hold things like this and pretend I&#39;m a giant.&quot; Line by Craig Bierko. Funniest man I have ever met. <a href="https://t.co/b1TmUuckjH">pic.twitter.com/b1TmUuckjH</a></p>&mdash; matthew perry (@MatthewPerry) <a href="https://twitter.com/MatthewPerry/status/1329015312004632576?ref_src=twsrc%5Etfw">November 18, 2020</a></blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>                        <div class="image">            <img src=" http://127.0.0.1:7001/dega/test-space/2020/10/1605767086916_test.png " id="test-png" alt="Test alt text">            <p>A Caption</p>            </div>                        <div class="image">        <img src=" http://127.0.0.1:7001/dega/test-space/2020/10/1605767086916_test.png " id="test-png" alt="Test alt text">        <p>A Caption</p>        </div>
	`

	t.Run("convert sample editorjs description to html", func(t *testing.T) {
		editorjsBlocks := make(map[string]interface{})
		err := json.Unmarshal(testDescription.RawMessage, &editorjsBlocks)
		if err != nil {
			t.Error(err)
		}

		viper.Set("templates_path", "templates")
		html, err := EditorjsToHTML(editorjsBlocks)
		if err != nil {
			t.Error(err)
		}

		t.Log("\n", html, "\n")

		outHtml = strings.TrimSpace(outHtml)
		outHtml = strings.ReplaceAll(outHtml, "\n", "")
		outHtml = strings.ReplaceAll(outHtml, "\t", "")

		if html != outHtml {
			t.Fail()
		}
	})

	t.Run("throwing error when new block is found", func(t *testing.T) {
		descString := `{
		"time": 1605767087876,
		"blocks": [
			{
				"type": "new",
				"data": {}
			}
		],
		"version": "2.19.0"
		}
		`
		testDescription := postgres.Jsonb{
			RawMessage: []byte(descString),
		}

		editorjsBlocks := make(map[string]interface{})
		err := json.Unmarshal(testDescription.RawMessage, &editorjsBlocks)
		if err != nil {
			t.Error(err)
		}

		viper.Set("templates_path", "templates")
		_, err = EditorjsToHTML(editorjsBlocks)

		if err == nil || err.Error() != "unparsed block found in description" {
			t.Fail()
		}

	})
}
