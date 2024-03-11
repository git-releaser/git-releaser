package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/config"
	"net/http"
	"strconv"
)

func (g Client) CreateRelease(baseBranch string, version config.Versions, description string) error {
	err := g.createTag(g.ProjectID, baseBranch, version, description)
	if err != nil {
		return err
	}

	if len(g.PropagationTargets) > 0 {
		fmt.Println("Propagating release to other repositories...")
		for _, target := range g.PropagationTargets {
			if target.TargetBranch == "" {
				target.TargetBranch = baseBranch
			}

			projectId, err := strconv.Atoi(target.Target)
			if err != nil {
				return err
			}

			err = g.createTag(projectId, target.TargetBranch, version, description)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g Client) CheckRelease(version config.Versions) (bool, error) {
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/repository/tags", g.ApiURL, g.ProjectID),
		Method: http.MethodGet,
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(req)
		return false, fmt.Errorf("failed to fetch tags. Status code: %d", resp.StatusCode)
	}

	var tags []config.Tag
	if err := json.Unmarshal(resp.Body, &tags); err != nil {
		fmt.Println("Test")
		return false, err
	}

	// Check if the desired tag is in the list
	for _, tag := range tags {
		if tag.Name == version.CurrentVersion.Original() {
			return true, nil
		}
	}
	return false, nil
}

func (g Client) createTag(project int, baseBranch string, version config.Versions, description string) error {
	var err error
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/releases", g.ApiURL, project),
		Method: http.MethodPost,
	}

	payload := map[string]interface{}{
		"tag_name":    version.CurrentVersion.Original(),
		"ref":         baseBranch,
		"description": description,
	}

	req.Payload, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	if g.DryRun {
		fmt.Println("Dry run: would create release with the following data:")
		fmt.Printf("Tag name: %s\n", version.CurrentVersion.Original())
		fmt.Printf("Ref: %s\n", baseBranch)
		fmt.Printf("Description: %s\n", description)
		return nil
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create release. Status code: %d", resp.StatusCode)
	}

	fmt.Println("Release created successfully (" + req.URL + ")")
	return nil
}
