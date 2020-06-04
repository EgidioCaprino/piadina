package godev

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func FetchResultHtmlPage(page int, query string) (string, error) {
	if page <= 0 {
		return "", errors.New("page parameter should be equal or greater to 1")
	}
	escapedQuery := url.QueryEscape(query)
	targetUrl := fmt.Sprintf("https://pkg.go.dev/search?page=%d&q=%s", page, escapedQuery)
	response, err := http.Get(targetUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response status %d: %s", response.StatusCode, response.Status)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(responseBytes), nil
}
