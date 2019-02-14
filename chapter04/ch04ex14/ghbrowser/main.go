// Ghbrowser runs an HTTP server to browse the issues of a GitHub repository.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bewuethr/gopl/chapter04/ch04ex14/github"
)

type headerInfo struct {
	Title string
	Index bool
}

var (
	issues     []github.Issue
	users      []github.User
	milestones []github.Milestone
)

func init() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: ghbrowser OWNER REPO")
		os.Exit(1)
	}
	var err error
	issues, users, milestones, err = github.GetRepoInfo(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/issue/", issueHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/milestone/", milestoneHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if err := header.Execute(w, headerInfo{"Issue list", true}); err != nil {
		log.Fatal(err)
	}
	if err := issueList.Execute(w, issues); err != nil {
		log.Fatal(err)
	}
}

func issueHandler(w http.ResponseWriter, r *http.Request) {
	pathWords := strings.Split(r.URL.Path, "/")
	n, err := strconv.Atoi(pathWords[len(pathWords)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	issue, err := github.GetIssueByNumber(issues, n)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := header.Execute(w, headerInfo{fmt.Sprintf("Issue %d", n), false}); err != nil {
		log.Fatal(err)
	}
	if err := issueTempl.Execute(w, issue); err != nil {
		log.Fatal(err)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	pathWords := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(pathWords[len(pathWords)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := github.GetUserByID(users, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := header.Execute(w, headerInfo{fmt.Sprintf("User %d", id), false}); err != nil {
		log.Fatal(err)
	}
	if err := userTempl.Execute(w, user); err != nil {
		log.Fatal(err)
	}
}

func milestoneHandler(w http.ResponseWriter, r *http.Request) {
	pathWords := strings.Split(r.URL.Path, "/")
	n, err := strconv.Atoi(pathWords[len(pathWords)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	milestone, err := github.GetMilestoneByNumber(milestones, n)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := header.Execute(w, headerInfo{fmt.Sprintf("Milestone %d", n), false}); err != nil {
		log.Fatal(err)
	}
	if err := milestoneTempl.Execute(w, milestone); err != nil {
		log.Fatal(err)
	}
}
