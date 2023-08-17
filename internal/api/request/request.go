package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Eg: http://192.168.0.1/rest/v1/
const UrlAPI = "%s://%s/rest/%s/%s"

type Request[T any] struct {
	URL      url.URL
	Response *T
}

type CollectionResult struct {
	TotalElementsCount    int `json:"total_elements_count"`
	FilteredElementsCount int `json:"filtered_elements_count"`
}

func New[T any](host, path string, response T) *Request[T] {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   fmt.Sprintf("/rest/%s/%s", "v8", path),
	}

	return &Request[T]{
		URL:      u,
		Response: &response,
	}
}

func (r *Request[T]) GetPretty() (string, error) {
	resBody, err := r.Get()
	if err != nil {
		return "", fmt.Errorf("can't read response body: %w", err)
	}

	pretty, err := jsonIndent(resBody)
	if err != nil {
		return "", err
	}

	return pretty, nil
}

func (r *Request[T]) Get() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, r.URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create get request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	return io.ReadAll(res.Body)
}

func jsonIndent(data []byte) (string, error) {
	pretty := bytes.Buffer{}
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return "", fmt.Errorf("can't pretty marshal data: %w", err)
	}

	return pretty.String(), nil
}

func (r *Request[T]) GetUnmarshalled() (*T, error) {
	body, err := r.Get()
	if err != nil {
		return r.Response, fmt.Errorf("could not complete API request: %w", err)
	}

	if err := json.Unmarshal(body, r.Response); err != nil {
		return r.Response, fmt.Errorf("can't unmarshal JSON response: %w", err)
	}

	return r.Response, nil
}
