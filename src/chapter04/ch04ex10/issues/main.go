// Issues prints a table of GitHub issues matching the search terms.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"chapter04/ch04ex10/github"
)

const (
	aMonth = 30 * 24 * time.Hour
	aYear  = 365 * 24 * time.Hour
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Indices of items in different age categories
	ageCats := make(map[string][]int)
	for i, item := range result.Items {
		switch {
		case time.Since(item.CreatedAt) > aYear:
			ageCats["year"] = append(ageCats["year"], i)
		case time.Since(item.CreatedAt) > aMonth:
			ageCats["month"] = append(ageCats["month"], i)
		default:
			ageCats["new"] = append(ageCats["new"], i)
		}
	}

	fmt.Printf("%d issues total, of which on first page...\n", result.TotalCount)

	fmt.Printf("%d issues more recent than one month:\n", len(ageCats["new"]))
	for _, idx := range ageCats["new"] {
		fmt.Printf("\t#%-5d %9.9s %.55s\n",
			result.Items[idx].Number, result.Items[idx].User.Login, result.Items[idx].Title)
	}

	fmt.Printf("%d issues older than one month:\n", len(ageCats["month"]))
	for _, idx := range ageCats["month"] {
		fmt.Printf("\t#%-5d %9.9s %.55s\n",
			result.Items[idx].Number, result.Items[idx].User.Login, result.Items[idx].Title)
	}

	fmt.Printf("%d issues older than one year:\n", len(ageCats["year"]))
	for _, idx := range ageCats["year"] {
		fmt.Printf("\t#%-5d %9.9s %.55s\n",
			result.Items[idx].Number, result.Items[idx].User.Login, result.Items[idx].Title)
	}
}
