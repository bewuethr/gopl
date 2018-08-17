package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRepoInfo gets all issues, users and milestones for the given repo
func GetRepoInfo(owner, repo string) ([]Issue, []User, []Milestone, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", apiBaseURL, owner, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	qString := req.URL.Query()
	qString.Add("direction", "asc")
	qString.Add("state", "all")
	req.URL.RawQuery = qString.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, nil, nil, err
	}
	users, milestones := extractUsersMilestones(issues)
	return issues, users, milestones, nil
}

func extractUsersMilestones(issues []Issue) ([]User, []Milestone) {
	uSeen, mSeen := make(map[int]bool), make(map[int]bool)
	var users []User
	var milestones []Milestone

	for _, issue := range issues {
		if !uSeen[issue.User.ID] {
			users = append(users, *issue.User)
			uSeen[issue.User.ID] = true
		}
		if issue.Milestone != nil && !mSeen[issue.Milestone.Number] {
			milestones = append(milestones, *issue.Milestone)
			mSeen[issue.Milestone.Number] = true
		}
	}

	return users, milestones
}

// GetIssueByNumber returns the issue with number n from the given slice of
// issues.
func GetIssueByNumber(issues []Issue, n int) (*Issue, error) {
	for _, issue := range issues {
		if issue.Number == n {
			return &issue, nil
		}
	}
	return nil, fmt.Errorf("could not find issue with number %d", n)
}

// GetUserByID returns the user with ID id from the given slice of users.
func GetUserByID(users []User, id int) (*User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("could not find user with ID %d", id)
}
