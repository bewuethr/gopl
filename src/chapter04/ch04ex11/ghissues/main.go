// ghissues lets you create, read, update and close GitHub issues from the
// command line.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"chapter04/ch04ex11/github"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "No arguments provided\n"+
			"Usage: ghissues <command> [<args>]\n"+
			"where <command> is \"read\", \"create\", \"update\" or \"close\"\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "read":
		handleRead()
	case "create":
		handleCreate()
	case "update":
		handleUpdate()
	case "close":
		handleClose()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command '%s'\n"+
			"Usage: ghissues <command> [<args>]\n"+
			"where <command> is \"read\", \"create\", \"update\" or \"close\"\n", os.Args[1])
		os.Exit(1)
	}
}

func handleRead() {
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
}

func printIssue(issue *github.Issue) {
	fmt.Printf("#%d - %s\n%s\n\n", issue.Number, issue.Title, issue.HTMLURL)
	fmt.Printf("Created: %s\nUser: %s\n\n", issue.CreatedAt, issue.User.Login)
	fmt.Printf("State: %s\n\n%s\n", issue.State, issue.Body)
}

func handleCreate() {
	if len(os.Args) < 5 {
		fmt.Fprint(os.Stderr, "Not enough arguments provided\n"+
			"Usage: ghissues create OWNER REPO TITLE [BODY]\n")
		os.Exit(1)
	}
	owner, repo, title := os.Args[2], os.Args[3], os.Args[4]
	var body string
	if len(os.Args) >= 6 {
		body = os.Args[5]
	} else {
		body = getFromEditor()
	}

	result, err := github.CreateIssue(owner, repo, title, body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Issue successfully created at %s\n", result.HTMLURL)
}

func handleUpdate() {
	if len(os.Args) < 6 {
		fmt.Fprint(os.Stderr, "Not enough arguments provided\n"+
			"Usage: ghissues update OWNER REPO NUMBER TITLE [BODY]\n")
		os.Exit(1)
	}
	owner, repo, title := os.Args[2], os.Args[3], os.Args[5]
	number, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	var body string
	if len(os.Args) >= 7 {
		body = os.Args[6]
	} else {
		body = getFromEditor()
	}

	result, err := github.UpdateIssue(owner, repo, number, title, body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Issue at %s successfully updated\n", result.HTMLURL)
}

func handleClose() {
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
}

func getFromEditor() string {
	var editor string
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		editor = "vim"
	}
	tmpfile, err := ioutil.TempFile("", "body.md")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Open preferred editor
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Get contents from tempfile
	body, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
