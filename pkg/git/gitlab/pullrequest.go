package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thschue/git-releaser/pkg/changelog"
	"github.com/thschue/git-releaser/pkg/naming"
	"io"
	"net/http"
)

type MergeRequest struct {
	ID           int    `json:"id"`
	IID          int    `json:"iid"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Title        string `json:"title"`
}

func (g Client) CheckCreatePullRequest(source string, target string, currentVersion string, nextVersion string) error {
	err := g.CreatePullRequest(source, target, currentVersion, nextVersion)
	if err != nil {
		return err
	}
	return nil
}

func (g Client) CreatePullRequest(source string, target string, currentVersion string, nextVersion string) error {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID)

	// Check if a pull request with the same source and target branches already exists
	existingPrID, err := g.getExistingPullRequestID(source, target)
	if err != nil {
		return err
	}

	commits, _ := g.GetCommitsSinceRelease(currentVersion)
	conventionalCommits := changelog.ParseConventionalCommits(commits)
	changelog := changelog.GenerateChangelog(conventionalCommits, g.ProjectURL)

	title := naming.GeneratePrTitle(nextVersion)
	description := naming.CreatePrDescription(nextVersion, changelog)

	payload := map[string]interface{}{
		"source_branch": source,
		"target_branch": target,
		"title":         title,
		"description":   description,
	}

	var req *http.Request

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if existingPrID != 0 {
		// If the pull request already exists, update its description
		url = fmt.Sprintf("%s/projects/%d/merge_requests/%d", g.ApiURL, g.ProjectID, existingPrID)
		req, err = http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
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
	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

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

	if existingPrID != 0 {
		fmt.Println("Pull request updated successfully.")
	} else {
		fmt.Println("Pull request created successfully.")
	}

	return nil
}

// getExistingPullRequestID retrieves the ID of an existing pull request with the same source and target branches
func (g Client) getExistingPullRequestID(source, target string) (int, error) {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID)

	// Fetch all merge requests
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to fetch merge requests. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	var mergeRequests []struct {
		ID           int    `json:"id"`
		IID          int    `json:"iid"`
		SourceBranch string `json:"source_branch"`
		TargetBranch string `json:"target_branch"`
		State        string `json:"state"`
		Title        string `json:"title"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mergeRequests); err != nil {
		return 0, err
	}

	// Find the ID of the existing pull request with the same source and target branches
	for _, pr := range mergeRequests {
		if pr.SourceBranch == source && pr.TargetBranch == target && pr.State == "opened" {
			if pr.IID != 0 {
				return pr.IID, nil
			}
			return pr.ID, nil
		}
	}

	return 0, nil // No existing pull request found
}
