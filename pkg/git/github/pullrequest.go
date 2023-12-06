package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thschue/git-releaser/pkg/changelog"
	"github.com/thschue/git-releaser/pkg/naming"
	"io"
	"net/http"
)

type PullRequest struct {
	ID           int    `json:"id"`
	Number       int    `json:"number"`
	Title        string `json:"title"`
	State        string `json:"state"`
	SourceBranch string `json:"head_ref"`
	TargetBranch string `json:"base_ref"`
}

func (g Client) CheckCreatePullRequest(source string, target string, currentVersion string, nextVersion string) error {
	err := g.CreatePullRequest(source, target, currentVersion, nextVersion)
	if err != nil {
		return err
	}
	return nil
}

func (g Client) CreatePullRequest(source string, target string, currentVersion string, nextVersion string) error {
	url := fmt.Sprintf("%s/repos/%s/pulls", g.ApiURL, g.Repository)

	// Check if a pull request with the same source and target branches already exists
	existingPrNumber, err := g.getExistingPullRequestNumber(source, target)
	if err != nil {
		return err
	}

	commits, _ := g.GetCommitsSinceRelease(currentVersion)
	conventionalCommits := changelog.ParseConventionalCommits(commits)
	changelog := changelog.GenerateChangelog(conventionalCommits, g.ProjectURL)

	title := naming.GeneratePrTitle(nextVersion)
	description := naming.CreatePrDescription(nextVersion, changelog)

	payload := map[string]interface{}{
		"title": title,
		"body":  description,
		"head":  source,
		"base":  target,
	}

	var req *http.Request

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if existingPrNumber != 0 {
		// If the pull request already exists, update its description
		url = fmt.Sprintf("%s/repos/%s/pulls/%d", g.ApiURL, g.Repository, existingPrNumber)
		req, err = http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return err
		}
	} else {
		// If the pull request doesn't exist, create a new one
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create/update pull request. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	if existingPrNumber != 0 {
		fmt.Println("Pull request updated successfully.")
	} else {
		fmt.Println("Pull request created successfully.")
	}

	return nil
}

// getExistingPullRequestNumber retrieves the number of an existing pull request with the same source and target branches
func (g Client) getExistingPullRequestNumber(source, target string) (int, error) {
	url := fmt.Sprintf("%s/repos/%s/pulls", g.ApiURL, g.Repository)

	// Fetch all pull requests
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to fetch pull requests. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	var pullRequests []PullRequest

	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return 0, err
	}

	// Find the number of the existing pull request with the same source and target branches
	for _, pr := range pullRequests {
		if pr.SourceBranch == source && pr.TargetBranch == target && pr.State == "open" {
			return pr.Number, nil
		}
	}

	return 0, nil // No existing pull request found
}
