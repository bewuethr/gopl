// Package github provides a Go API for the GitHub issue tracker.
package github

import "time"

const apiBaseURL = "https://api.github.com"

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	Milestone *Milestone
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	ID        int
	Login     string
	HTMLURL   string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
}

type Milestone struct {
	Number      int
	HTMLURL     string `json:"html_url"`
	State       string
	Title       string
	Description string // in Markdown format
}
