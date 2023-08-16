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

type Request struct {
	URL url.URL
}

func New(host, path string) *Request {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   fmt.Sprintf("/rest/%s/%s", "v8", path),
	}

	return &Request{
		URL: u,
	}
}

func (r *Request) Get() (string, error) {
	req, err := http.NewRequest(http.MethodGet, r.URL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("can't create get request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("can't read response body: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("can't read response body: %w", err)
	}

	pretty, err := jsonIndent(resBody)
	if err != nil {
		return "", err
	}

	return pretty, nil
}

func jsonIndent(data []byte) (string, error) {
	pretty := bytes.Buffer{}
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return "", fmt.Errorf("can't pretty marshal data: %w", err)
	}

	return pretty.String(), nil
}
