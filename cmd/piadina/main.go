package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	godev "github.com/EgidioCaprino/piadina/internal"
)

const pageLimit = 10

func main() {
	flag.Parse()
	query := flag.Arg(0)
	if query == "" {
		log.Fatalln("expected a query")
	}
	results, err := godev.QueryGoDev(query, pageLimit)
	if err != nil {
		log.Fatalln(err)
	}
	sortedResults, err := sortResults(results)
	if err != nil {
		log.Fatalln(err)
	}
	for _, resultPackage := range sortedResults {
		fmt.Println(resultPackage)
		fmt.Println("----------")
	}
}

func sortResults(results []godev.ResultPackage) ([]godev.ResultPackage, error) {
	sorted := make([]godev.ResultPackage, len(results))
	copy(sorted, results)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Imports > sorted[j].Imports
	})
	return sorted, nil
}
