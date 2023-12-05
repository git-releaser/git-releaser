package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (g Client) CheckCreateBranch(version string) (string, error) {
	branchName := fmt.Sprintf("release-%s", version)
	branchExists, _ := g.branchExists(branchName)
	if !branchExists {
		err := g.createBranch(branchName)
		if err != nil {
			return "", err
		}
	}
	return branchName, nil
}

func (g Client) branchExists(branchName string) (bool, error) {
	url := fmt.Sprintf("%s/repos/%s/branches/%s", g.ApiURL, g.Repository, branchName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Branch exists
		return true, nil
	case http.StatusNotFound:
		// Branch does not exist
		return false, nil
	default:
		// Other status codes indicate an error
		return false, fmt.Errorf("failed to check branch existence. Status code: %d", resp.StatusCode)
	}
}

// createBranch creates a branch in a GitHub repository
func (g Client) createBranch(branchName string) error {
	url := fmt.Sprintf("%s/repos/%s/git/refs", g.ApiURL, g.Repository)
	baseSha, err := g.getBaseBranchSHA("main")
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"ref": fmt.Sprintf("refs/heads/%s", branchName),
		"sha": baseSha, // Replace with the commit SHA or branch you want to base the new branch on
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
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create branch. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	fmt.Printf("Branch '%s' created successfully.\n", branchName)
	return nil
}

func (g Client) getBaseBranchSHA(baseBranch string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/branches/%s", g.ApiURL, g.Repository, baseBranch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+g.AccessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch branch details. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var branchInfo map[string]interface{}
	if err := json.Unmarshal(body, &branchInfo); err != nil {
		return "", err
	}

	commit, ok := branchInfo["commit"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("commit information not found in the response")
	}

	sha, ok := commit["sha"].(string)
	if !ok {
		return "", fmt.Errorf("SHA not found in the commit information")
	}

	return sha, nil
}
