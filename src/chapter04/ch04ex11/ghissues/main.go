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

// TODO Read arguments with command line flags
// TODO Open editor if title/body is missing in arguments

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No arguments provided\n"+
			"Usage: ghissues <command> [<args>]\n"+
			"where <command> is \"read\", \"create\" or \"update\"\n")
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

	case "create":
		if len(os.Args) < 6 {
			fmt.Fprint(os.Stderr, "Not enough arguments provided\n"+
				"Usage: ghissues create OWNER REPO TITLE BODY\n")
			os.Exit(1)
		}
		owner, repo, title, body := os.Args[2], os.Args[3], os.Args[4], os.Args[5]
		result, err := github.CreateIssue(owner, repo, title, body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Issue successfully created at %s\n", result.HTMLURL)

	case "update":
		if len(os.Args) < 7 {
			fmt.Fprint(os.Stderr, "Not enough arguments provided\n"+
				"Usage: ghissues update OWNER REPO NUMBER TITLE BODY\n")
			os.Exit(1)
		}
		owner, repo, title, body := os.Args[2], os.Args[3], os.Args[5], os.Args[6]
		number, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal(err)
		}
		result, err := github.UpdateIssue(owner, repo, number, title, body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Issue at %s successfully updated\n", result.HTMLURL)

	case "close":
		if len(os.Args) < 5 {
			fmt.Fprint(os.Stderr, "Not enough arguments provided\n"+
				"Usage: ghissues close OWNER REPO NUMBER\n")
			os.Exit(1)
		}
		owner, repo := os.Args[2], os.Args[3]
		number, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatal(err)
		}
		result, err := github.CloseIssue(owner, repo, number)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Issue at %s successfully closed\n", result.HTMLURL)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command '%s'\n"+
			"Usage: ghissues <command> [<args>]\n"+
			"where <command> is \"read\", \"create\", \"update\" or \"close\"\n", os.Args[1])
		os.Exit(1)
	}
}

func printIssue(issue *github.Issue) {
	fmt.Printf("#%d - %s\n%s\n\n", issue.Number, issue.Title, issue.HTMLURL)
	fmt.Printf("Created: %s\nUser: %s\n\n", issue.CreatedAt, issue.User.Login)
	fmt.Printf("State: %s\n\n%s\n", issue.State, issue.Body)
}
