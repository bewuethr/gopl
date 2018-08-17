// Ghbrowser runs an HTTP server to browse the issues of a GitHub repository.
package main

// TODO Improve navigation and linking between templates
// TODO Maybe move templates into separate file
// TODO Some simple styling
// TODO Markdown renderer? For example, https://github.com/russross/blackfriday

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"chapter04/ch04ex14/github"
)

var issueList = template.Must(template.New("issuelist").Parse(`<h1>{{len .}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .}}
<tr>
  <td><a href='issue/{{.Number}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='user/{{.User.ID}}'>{{.User.Login}}</a></td>
  <td><a href='issue/{{.Number}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>`))

var issueTempl = template.Must(template.New("issue").Parse(`<h1>#{{.Number}}: {{.Title}}</h1>
<p>by <a href="/user/{{.User.ID}}">{{.User.Login}}</a> &ndash; {{.State}}</p>
{{if .Milestone}}<p>Milestone: {{.Milestone.Title}}</p>{{end}}
<p>Created at {{.CreatedAt}}</p>
<pre>{{.Body}}</pre>`))

var userTempl = template.Must(template.New("user").Parse(`<h1><a href="{{.HTMLURL}}">{{.Login}}</a></h1>
<img src="{{.AvatarURL}}" alt="Avatar"/>`))

var milestoneTempl = template.Must(template.New("milestone").Parse(`<h1>Milestone #{{.Number}}: {{.Title}}</h1>
<p><a href="{{.HTMLURL}}">Link</a> &ndash; {{.State}}</p>
<pre>{{.Description}}</pre>`))

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
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := userTempl.Execute(w, user); err != nil {
		log.Fatal(err)
	}
}
