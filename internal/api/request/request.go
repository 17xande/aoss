package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Eg: http://192.168.0.1/rest/v1/
const UrlAPI = "%s://%s/rest/%s/%s"

type loginResponse struct {
	PayloadSize int
	Uri         string
	Cookie      string
}

type CollectionResult struct {
	TotalElementsCount    int `json:"total_elements_count"`
	FilteredElementsCount int `json:"filtered_elements_count"`
}

func GetUnmarshalled[T any](host, path string, cookie *http.Cookie, result *T) error {
	res, err := get(host, path, cookie)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error unmarshalling result: %w", err)
	}

	return nil
}

func get(host, path string, cookie *http.Cookie) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/rest/v8/%s", host, path)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if cookie != nil {
		r.AddCookie(cookie)
	}

	return http.DefaultClient.Do(r)
}

func post(host, path string, cookie *http.Cookie, body []byte) (*http.Response, error) {

	return nil, nil
}

func GetJson(host, path string, cookie *http.Cookie) (string, error) {
	res, err := get(host, path, cookie)
	if err != nil {
		return "", err
	}

	jsonRes, err := json.MarshalIndent(res.Body, "", "  ")
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return string(jsonRes), fmt.Errorf("status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return string(jsonRes), err
	}

	indented := bytes.Buffer{}
	if err := json.Indent(&indented, body, "", "  "); err != nil {
		return string(jsonRes), err
	}

	return indented.String(), nil
}

func authenticate(host string) (*http.Cookie, error) {
	creds, err := readCreds()
	if err != nil {
		return nil, err
	}

	path := "/login-sessions"
	res, err := post(host, path, nil, creds)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var loginRes loginResponse

	if err := json.Unmarshal(body, &loginRes); err != nil {
		return nil, err
	}

	token := strings.Split(loginRes.Cookie, "=")[1]

	cookie := &http.Cookie{
		Name:   "sessionId",
		Value:  token,
		MaxAge: 90,
	}

	return cookie, nil
}

func readCreds() ([]byte, error) {
	creds, err := os.ReadFile("creds.json")
	if err != nil {
		return nil, fmt.Errorf("can't read credential file: %w", err)
	}

	return creds, nil
}
