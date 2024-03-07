package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/naming"
	"io"
	"net/http"
)

func (g Client) CheckCreatePullRequest(source string, target string, versions config.Versions) error {
	err := g.createPullRequest(source, target, versions)
	if err != nil {
		return err
	}

	// Check if other git-releaser pull requests exist and close them
	err = g.closeOldPullRequests(source)
	if err != nil {
		return err
	}

	return nil
}

func (g Client) createPullRequest(source string, target string, versions config.Versions) error {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID)

	// Check if a pull request with the same source and target branches already exists
	existingPrID, err := g.getExistingPullRequestID(source, target)
	if err != nil {
		return err
	}

	commits, _ := g.GetCommitsSinceRelease(versions.CurrentVersion.Original())
	conventionalCommits := changelog.ParseConventionalCommits(commits)
	changelog := changelog.GenerateChangelog(conventionalCommits, g.ProjectURL)

	title := naming.GeneratePrTitle(versions.NextVersion.Original())
	description := naming.CreatePrDescription(versions.NextVersion.Original(), changelog)

	payload := map[string]interface{}{
		"source_branch": source,
		"target_branch": target,
		"title":         title,
		"description":   description,
		"labels":        []string{"release"},
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

	if g.DryRun {
		fmt.Println("Dry run: pull request would be created with the following details:")
		fmt.Println("Title: " + title)
		fmt.Println("Description: " + description)
		fmt.Println("Source branch: " + source)
		fmt.Println("Target branch: " + target)
		return nil
	}

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

func (g Client) closeOldPullRequests(currentSource string) error {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID)

	// Fetch all merge requests
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch merge requests. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	var mergeRequests []MergeRequest

	if err := json.NewDecoder(resp.Body).Decode(&mergeRequests); err != nil {
		return err
	}

	for _, mr := range mergeRequests {
		// Check if the merge request is open and has a "release" label
		if mr.State == "opened" && contains(mr.Labels, "release") && mr.SourceBranch != currentSource {
			// Close the merge request
			err := g.closeMergeRequest(mr.IID)
			if err != nil {
				return err
			}

			// Delete the source branch
			err = g.deleteBranch(mr.SourceBranch)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g Client) closeMergeRequest(id int) error {
	url := fmt.Sprintf("%s/projects/%d/merge_requests/%d", g.ApiURL, g.ProjectID, id)

	payload := map[string]interface{}{
		"state_event": "close",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to close merge request. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
