package requestx

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Request does a http request
func Request(method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
