package godev

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseResultHtmlPage(t *testing.T) {
	inputFile, err := filepath.Abs(filepath.Join(".", "example_result_page.html"))
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	resultPackages, err := parseResultHtmlPage(file)
	if err != nil {
		t.Fatal(err)
	}

	expectedResults := 10
	if len(resultPackages) != expectedResults {
		t.Fatalf("expected to get %d some packages, got %d", expectedResults, len(resultPackages))
	}

	expectedFirst := ResultPackage{
		Name:        "log",
		Description: "Package log implements a simple logging package.",
		Version:     "go1.14.4",
		PublishDate: "5 days ago",
		Imports:     150652,
		License:     "BSD-3-Clause",
	}
	if resultPackages[0] != expectedFirst {
		t.Errorf("expected first package to be %+v but got %+v", expectedFirst, resultPackages[0])
	}
}

func TestParseInfo(t *testing.T) {
	info := `Version: go1.14.4
		|
		Published: 3 days ago
		|
		Imported by: 149996
		|
		License:
		BSD-3-Clause`
	resultPackage, err := parseInfo(info)
	if err != nil {
		t.Fatal(err)
	}

	expectedVersion := "go1.14.4"
	if resultPackage.Version != expectedVersion {
		t.Errorf("expected version to be %s but got %s", expectedVersion, resultPackage.Version)
	}

	expectedPublishDate := "3 days ago"
	if resultPackage.PublishDate != expectedPublishDate {
		t.Errorf("expected version to be %s but got %s", expectedPublishDate, resultPackage.PublishDate)
	}

	expectedImports := 149996
	if resultPackage.Imports != expectedImports {
		t.Errorf("expected imports to be %d but got %d", expectedImports, resultPackage.Imports)
	}

	expectedLicense := "BSD-3-Clause"
	if resultPackage.License != expectedLicense {
		t.Errorf("expected version to be %s but got %s", expectedLicense, resultPackage.License)
	}
}
