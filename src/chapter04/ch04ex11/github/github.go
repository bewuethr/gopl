// Package github provides a Go API for the GitHub issue tracker.
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const BaseAPIURL = "https://api.github.com"

var authToken string

func init() {
	var ok bool
	authToken, ok = os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("Environment variable GITHUB_TOKEN is not set")
	}
}

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

type reqBody struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"` // in Markdown format
	State string `json:"state,omitempty"`
}

type URLResponse struct {
	HTMLURL string `json:"html_url"`
}

// ReadIssue gets an issue from a repo.
func ReadIssue(owner, repo string, number int) (*Issue, error) {
	r, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/repos/%s/%s/issues/%d", BaseAPIURL, owner, repo, number), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Accept", "application/vnd.github.v3.+json")
	r.Header.Set("User-Agent", "gopl-ch04ex11-cli-tool")
	resp, err := http.DefaultClient.Do(r)

	// We must close resp.Body on all execution paths.
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

// CreateIssue creates a new issue in the specified repo and returns the
// response struct with the URL if successful.
func CreateIssue(owner, repo, title, body string) (*URLResponse, error) {
	rBody := reqBody{
		Title: title,
		Body:  body,
	}
	url := fmt.Sprintf("%s/repos/%s/%s/issues", BaseAPIURL, owner, repo)
	resp, err := sendReqWithBody(rBody, url, http.MethodPost, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateIssue updates an issue and returns the response struct with the URL if
// successful.
func UpdateIssue(owner, repo string, number int, title, body string) (*URLResponse, error) {
	rBody := reqBody{
		Title: title,
		Body:  body,
	}
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", BaseAPIURL, owner, repo, number)
	resp, err := sendReqWithBody(rBody, url, http.MethodPatch, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CloseIssue closes an issue and returns the response struct with the URL if
// successful.
func CloseIssue(owner, repo string, number int) (*URLResponse, error) {
	rBody := reqBody{
		State: "closed",
	}
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", BaseAPIURL, owner, repo, number)
	resp, err := sendReqWithBody(rBody, url, http.MethodPatch, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// sendReqWithBody is what the create, update and close requests have in common.
func sendReqWithBody(rBody reqBody, url, method string, expectStatus int) (*URLResponse, error) {
	jsonBody, err := json.Marshal(rBody)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	r.Header.Set("Accept", "application/vnd.github.v3.+json")
	r.Header.Set("User-Agent", "gopl-ch04ex11-cli-tool")
	r.Header.Set("Authorization", "Token "+authToken)
	resp, err := http.DefaultClient.Do(r)

	if resp.StatusCode != expectStatus {
		resp.Body.Close()
		return nil, fmt.Errorf("query failed: %s", resp.Status)
	}

	var response URLResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &response, nil
}
