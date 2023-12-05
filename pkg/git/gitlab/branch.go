package gitlab

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
	url := fmt.Sprintf("%s/projects/%d/repository/branches/%s", g.ApiURL, g.ProjectID, branchName)

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

func (g Client) createBranch(branchName string) error {
	url := fmt.Sprintf("%s/projects/%d/repository/branches", g.ApiURL, g.ProjectID)

	payload := map[string]interface{}{
		"branch": branchName,
		"ref":    "main",
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
		return fmt.Errorf("failed to create branch. Status code: %d, Body: %s", resp.StatusCode, body)
	}

	fmt.Printf("Branch '%s' created successfully.\n", branchName)
	return nil
}
