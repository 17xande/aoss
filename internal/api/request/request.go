package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/17xande/configdir"
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

type Auth struct {
	Host string
	http.Cookie
}

func GetUnmarshalled[T any](host, path string, auth *Auth, result *T) error {
	res, err := get(host, path, &auth.Cookie)
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

func PostUnmarshalled[T any](host, path string, auth *Auth, req *T) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("can't marshall request: %w", err)
	}

	reader := bytes.NewReader(body)

	res, err := post(host, path, &auth.Cookie, reader)
	if err != nil {
		return nil, fmt.Errorf("can't run PUT request: %w", err)
	}

	return res, nil
}

func put(host, path string, cookie *http.Cookie, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/rest/v8/%s", host, path)

	r, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	if cookie != nil {
		r.AddCookie(cookie)
	}

	return http.DefaultClient.Do(r)
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

func post(host, path string, cookie *http.Cookie, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/rest/v8/%s", host, path)

	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("could not create post request: %w", err)
	}

	if cookie != nil {
		r.AddCookie(cookie)
	}

	return http.DefaultClient.Do(r)
}

func delete(host, path string, cookie *http.Cookie) (*http.Response, error) {
	url := fmt.Sprintf("http://%s/rest/v8/%s", host, path)

	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create delete request: %w", err)
	}

	if cookie != nil {
		r.AddCookie(cookie)
	}

	return http.DefaultClient.Do(r)
}

func GetJson(host, path string) (string, error) {
	auth := Auth{
		Host: host,
	}

	if err := auth.Login(); err != nil {
		return "", fmt.Errorf("could not authenticate: %w", err)
	}
	defer auth.Logout()

	res, err := get(host, path, &auth.Cookie)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return string(body), fmt.Errorf("status code: %d", res.StatusCode)
	}

	indented := bytes.Buffer{}
	if err := json.Indent(&indented, body, "", "  "); err != nil {
		return string(body), err
	}

	return indented.String(), nil
}

func (a *Auth) Login() error {
	creds, err := readCreds()
	if err != nil {
		return err
	}

	reader := bytes.NewReader(creds)

	path := "/login-sessions"
	res, err := post(a.Host, path, nil, reader)
	if err != nil {
		return fmt.Errorf("post failed: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("authentication failed with status: %d\n%s", res.StatusCode, body)
	}

	var loginRes loginResponse

	if err := json.Unmarshal(body, &loginRes); err != nil {
		return err
	}

	token := strings.Split(loginRes.Cookie, "=")[1]

	a.Cookie = http.Cookie{
		Name:   "sessionId",
		Value:  token,
		MaxAge: 90,
	}

	return nil
}

func (a *Auth) Logout() {
	path := "/login-sessions"
	res, err := delete(a.Host, path, &a.Cookie)
	if err != nil {
		fmt.Printf("error trying to logout: %v\n", err)
	}

	if res.StatusCode != http.StatusNoContent {
		fmt.Printf("logout failed. Status: %d; Body: %v", res.StatusCode, res)
	}
}

func readCreds() ([]byte, error) {
	path := configdir.LocalConfig("aoss/creds.json")
	creds, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read credential file: %w", err)
	}

	return creds, nil
}
