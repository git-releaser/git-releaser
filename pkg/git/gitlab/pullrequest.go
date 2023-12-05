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

func (g Client) CheckCreatePullRequest(source string, target string) error {
	err := g.CreatePullRequest(source, target)
	if err != nil {
		return err
	}
	return nil
}

func (g Client) CreatePullRequest(source string, target string) error {
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
		if resp.StatusCode == http.StatusConflict {
			fmt.Println("Pull request already exists.")
			return nil
		}
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create pull request. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	fmt.Println("Pull request created successfully.")
	return nil
}

func generatePRTitleAndDescription(version string) (string, string) {
	title := fmt.Sprintf("Release %s", version)
	description := fmt.Sprintf("This is a description for the new pull request for version %s.", version)
	return title, description
}
