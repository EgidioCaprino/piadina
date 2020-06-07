package godev

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

type ResultPackage struct {
	Name        string
	Description string
	Version     string
	PublishDate string // @todo should be date?
	Imports     int
	License     string
}

func (resultPackage ResultPackage) String() string {
	title := fmt.Sprintf("%s - %s", resultPackage.Name, resultPackage.Version)
	packageUrl := fmt.Sprintf("https://pkg.go.dev/%s", resultPackage.Name)
	imports := printer.Sprintf("Imports: %d", resultPackage.Imports)
	lines := []string{
		title,
		packageUrl,
		imports,
		resultPackage.Description,
		resultPackage.License,
		resultPackage.PublishDate,
	}
	return strings.Join(lines, "\n")
}

func QueryGoDev(query string, pageLimit int) ([]ResultPackage, error) {
	if pageLimit < 1 {
		return nil, fmt.Errorf("expected page limit to be at least 1 but got %d", pageLimit)
	}
	var results []ResultPackage
	for page := 1; page <= pageLimit; page++ {
		htmlPageReader, err := fetchResultHtmlPage(page, query)
		if err != nil {
			return nil, err
		}
		defer htmlPageReader.Close()

		pageResults, err := parseResultHtmlPage(htmlPageReader)
		if err != nil {
			return nil, err
		}
		if len(pageResults) == 0 {
			break
		}

		results = append(results, pageResults...)
	}
	return results, nil
}

func fetchResultHtmlPage(page int, query string) (io.ReadCloser, error) {
	if page <= 0 {
		return nil, errors.New("page parameter should be equal or greater to 1")
	}
	escapedQuery := url.QueryEscape(query)
	targetUrl := fmt.Sprintf("https://pkg.go.dev/search?page=%d&q=%s", page, escapedQuery)
	response, err := http.Get(targetUrl)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status %d: %s", response.StatusCode, response.Status)
	}
	return response.Body, nil
}

func parseResultHtmlPage(htmlPage io.Reader) ([]ResultPackage, error) {
	document, err := goquery.NewDocumentFromReader(htmlPage)
	if err != nil {
		return nil, err
	}
	result := make([]ResultPackage, 0, 10)
	var firstError error
	document.Find("div.SearchSnippet").Each(func(index int, selection *goquery.Selection) {
		if firstError == nil {
			info := selection.Find("div.SearchSnippet-infoLabel").Text()
			resultPackage, err := parseInfo(info)
			if err != nil {
				firstError = fmt.Errorf("unable to parse info from %s: %w", info, err)
				return
			}
			resultPackage.Name = strings.TrimSpace(selection.Find("h2.SearchSnippet-header a").Text())
			resultPackage.Description = strings.TrimSpace(selection.Find("p.SearchSnippet-synopsis").Text())
			result = append(result, resultPackage)
		}
	})
	return result, firstError
}

func parseInfo(info string) (ResultPackage, error) {
	result := ResultPackage{}
	sections := strings.Split(info, "|")
	for _, section := range sections {
		parts := strings.Split(section, ":")
		title := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch title {
		case "Version":
			result.Version = value
			break
		case "Published":
			result.PublishDate = value
			break
		case "Imported by":
			imports, err := strconv.Atoi(value)
			if err != nil {
				return result, fmt.Errorf("unable to parse version from %s: %w", value, err)
			}
			result.Imports = imports
			break
		case "License":
			result.License = value
			break
		}
	}
	return result, nil
}
