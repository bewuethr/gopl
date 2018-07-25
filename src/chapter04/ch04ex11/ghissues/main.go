// ghissues lets you create, read, update and close GitHub issues from the
// command line.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"chapter04/ch04ex11/github"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No arguments provided\n"+
			"Usage: ghissues read OWNER REPO NUMBER\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "read":
		if len(os.Args) < 5 {
			fmt.Fprint(os.Stderr, "Not enough arguments provided\n"+
				"Usage: ghissues read OWNER REPO NUMBER\n")
			os.Exit(1)
		}
		owner, repo := os.Args[2], os.Args[3]
		number, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal(err)
		}
		result, err := github.ReadIssue(owner, repo, number)
		if err != nil {
			log.Fatal(err)
		}
		printIssue(result)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command '%s'\n"+
			"Usage: ghissues read OWNER REPO NUMBER\n", os.Args[1])
		os.Exit(1)
	}
}

func printIssue(issue *github.Issue) {
	fmt.Printf("%+v\n", *issue)
}
