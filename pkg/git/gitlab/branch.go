package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/naming"
	"net/http"
)

func (g Client) CheckCreateBranch(baseBranch string, version string, prefix string) (string, error) {
	branchName := naming.CreateBranchName(prefix, version)

	branchExists, _ := g.branchExists(branchName)
	if !branchExists {
		err := g.createBranch(baseBranch, branchName)
		if err != nil {
			return "", err
		}
	}
	return branchName, nil
}
func (g Client) branchExists(branchName string) (bool, error) {
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/repository/branches/%s", g.ApiURL, g.ProjectID, branchName),
		Method: http.MethodGet,
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return false, err
	}

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

func (g Client) createBranch(baseBranch string, branchName string) error {
	var err error
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/repository/branches", g.ApiURL, g.ProjectID),
		Method: http.MethodPost,
	}

	payload := map[string]interface{}{
		"branch": branchName,
		"ref":    baseBranch,
	}

	req.Payload, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	if g.DryRun {
		fmt.Printf("Dry run: Branch '%s' would be created.\n", branchName)
		return nil
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create branch. Status code: %d, Body: %s", resp.StatusCode, resp.Body)
	}

	fmt.Printf("Branch '%s' created successfully.\n", branchName)
	return nil
}

func (g Client) deleteBranch(branchName string) error {
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/repository/branches/%s", g.ApiURL, g.ProjectID, branchName),
		Method: http.MethodDelete,
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete branch. Status code: %d, Body: %s", resp.StatusCode, resp.Body)
	}

	fmt.Printf("Branch '%s' deleted successfully.\n", branchName)
	return nil
}
