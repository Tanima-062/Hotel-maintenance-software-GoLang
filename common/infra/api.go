package infra

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
)

// APIClient 外部APIのリクエスト設定用
type APIClient struct {
	URL      string
	Data     interface{}
	Response interface{}
}

// Post json post
func (a *APIClient) Post() error {

	postData, parseErr := json.Marshal(a.Data)
	if parseErr != nil {
		return parseErr
	}

	resp, postErr := http.Post(a.URL, "application/json", bytes.NewReader(postData))
	if postErr != nil {
		return postErr
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body

	// debug用。ログに出力されるようになります
	// body = io.TeeReader(body, os.Stderr)

	return json.NewDecoder(body).Decode(&a.Response)
}

// PostXml xml post
func (a *APIClient) PostXml(data []byte) error {
	resp, postErr := http.Post(a.URL, "text/xml", bytes.NewReader(data))
	if postErr != nil {
		return postErr
	}
	defer resp.Body.Close()
	byteValue, _ := ioutil.ReadAll(resp.Body)

	return xml.Unmarshal(byteValue, &a.Response)
}

// Get json get
func (a *APIClient) Get() error {
	resp, postErr := http.Get(a.URL)
	if postErr != nil {
		return postErr
	}
	defer resp.Body.Close()

	body, loadErr := ioutil.ReadAll(resp.Body)
	if loadErr != nil {
		return loadErr
	}
	return json.Unmarshal(body, &a.Response)
}
