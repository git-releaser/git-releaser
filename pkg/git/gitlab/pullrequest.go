package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thschue/git-releaser/pkg/versioning"
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

func (g GitLabClient) CheckCreatePullRequest(source string, target string) error {
	// check if branch exists in gitlab
	/* exists, err := g.doesPullRequestExist(source)
	if err != nil {
		return err
	}
	*/

	err := g.CreatePullRequest(source, target)
	if err != nil {
		return err
	}

	return nil
}

func (g GitLabClient) doesPullRequestExist(sourceBranch string) (bool, error) {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to get merge requests. Status code: %d", resp.StatusCode)
	}

	var mergeRequests []MergeRequest
	err = json.NewDecoder(resp.Body).Decode(&mergeRequests)
	if err != nil {
		return false, err
	}

	for _, mr := range mergeRequests {
		if mr.SourceBranch == sourceBranch && mr.TargetBranch == "main" {
			return true, nil
		}
	}

	return false, nil
}

func (g GitLabClient) CreatePullRequest(source string, target string) error {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID)

	version, _ := versioning.GetNextVersion()
	title, description := generatePRTitleAndDescription(version.String())
	payload := map[string]interface{}{
		"source_branch": source,
		"target_branch": target,
		"title":         title,
		"description":   description,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
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

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create pull request. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	fmt.Println("Pull request created successfully.")
	return nil
}

func checkIfBranchExists() {
	// check if branch exists in gitlab
}

func generatePRTitleAndDescription(version string) (string, string) {
	title := fmt.Sprintf("Release %s", version)
	description := fmt.Sprintf("This is a description for the new pull request for version %s.", version)
	return title, description
}
