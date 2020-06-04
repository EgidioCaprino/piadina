package godev

import (
	"testing"
)

func TestFetchResultHtmlPage(t *testing.T) {
	html, err := FetchResultHtmlPage(1, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(html) == 0 {
		t.Fatal("empty result string")
	}
}
