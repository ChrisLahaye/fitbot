package fitforfree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// API contains the internal state
type API struct {
	baseURL string
	headers map[string]string
}

// New initializes a new API instance
func New() *API {
	return &API{
		baseURL: "https://electrolyte.fitforfree.nl",
		headers: map[string]string{
			"App-Version": "4.6.3",
			"User-Agent":  "okhttp/4.2.0",
		},
	}
}

func (api *API) SetAuth(sessionID string) {
	api.headers["Authorization"] = fmt.Sprintf("Bearer %s", sessionID)
}

// Request makes an API request
func (api *API) Request(method string, path string, query interface{}, in interface{}, out interface{}) error {
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(data)
	}

	url, err := api.url(path, query)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	for key, value := range api.headers {
		req.Header.Add(key, value)
	}
	if in != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: %s", resp.Status, data)
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}

func (api *API) url(path string, query interface{}) (string, error) {
	url, err := url.Parse(api.baseURL)
	if err != nil {
		return "", err
	}

	url.Path = path
	url.RawQuery, err = querystring(query)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func querystring(v interface{}) (string, error) {
	values, err := query.Values(v)
	if err != nil {
		return "", err
	}

	values.Add("language", "en-US")
	return values.Encode(), nil
}
