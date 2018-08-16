package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetIssues gets all issues for the given repo
func GetIssues(owner, repo string) ([]Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	qString := req.URL.Query()
	qString.Add("direction", "asc")
	qString.Add("state", "all")
	req.URL.RawQuery = qString.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result []Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetIssueByNumber return the issue with number n from the given list of
// issues.
func GetIssueByNumber(issues []Issue, n int) (*Issue, error) {
	for _, issue := range issues {
		if issue.Number == n {
			return &issue, nil
		}
	}
	return nil, fmt.Errorf("could not find issue with number %d", n)
}
