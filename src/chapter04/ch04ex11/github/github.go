// Package github provides a Go API for the GitHub issue tracker.
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseAPIURL = "https://api.github.com"

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// ReadIssue gets an issue from a repo.
func ReadIssue(owner, repo string, number int) (*Issue, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/repos/%s/%s/issues/%d", BaseAPIURL, owner, repo, number), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.+json")
	req.Header.Set("User-Agent", "gopl-ch04ex11-cli-tool")
	resp, err := http.DefaultClient.Do(req)

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
