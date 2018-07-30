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

type CreateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"` // in Markdown format
}

type CreateResponse struct {
	URL string
}

type CloseRequest struct {
	State string `json:"state"`
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
func CreateIssue(owner, repo, title, body string) (*CreateResponse, error) {
	reqBody := CreateRequest{
		Title: title,
		Body:  body,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/repos/%s/%s/issues", BaseAPIURL, owner, repo),
		bytes.NewBuffer(jsonBody))
	req.Header.Set("Accept", "application/vnd.github.v3.+json")
	req.Header.Set("User-Agent", "gopl-ch04ex11-cli-tool")
	req.Header.Set("Authorization", "Token "+authToken)
	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return nil, fmt.Errorf("create query failed: %s", resp.Status)
	}

	var createResponse CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &createResponse, nil
}

// UpdateIssue updates an issue and returns the response struct with the URL if
// successful.
func UpdateIssue(owner, repo string, number int, title, body string) (*CreateResponse, error) {
	reqBody := CreateRequest{
		Title: title,
		Body:  body,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("%s/repos/%s/%s/issues/%d", BaseAPIURL, owner, repo, number),
		bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.+json")
	req.Header.Set("User-Agent", "gopl-ch04ex11-cli-tool")
	req.Header.Set("Authorization", "Token "+authToken)
	resp, err := http.DefaultClient.Do(req)

	// We must close resp.Body on all execution paths.
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

// CloseIssue closes an issue and returns the response struct with the URL if
// successful.
func CloseIssue(owner, repo string, number int) (*CreateResponse, error) {
	reqBody := CloseRequest{
		State: "closed",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("%s/repos/%s/%s/issues/%d", BaseAPIURL, owner, repo, number),
		bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.+json")
	req.Header.Set("User-Agent", "gopl-ch04ex11-cli-tool")
	req.Header.Set("Authorization", "Token "+authToken)
	resp, err := http.DefaultClient.Do(req)

	// We must close resp.Body on all execution paths.
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
