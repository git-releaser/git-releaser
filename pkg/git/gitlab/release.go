package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thschue/git-releaser/pkg/config"
	"net/http"
)

type Release struct {
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

type Tag struct {
	Name   string `json:"name"`
	Commit struct {
		ID string `json:"id"`
	} `json:"commit"`
}

func (g Client) CreateRelease(baseBranch string, version config.Versions, description string) error {
	url := fmt.Sprintf("%s/projects/%d/releases", g.ApiURL, g.ProjectID)

	payload := map[string]interface{}{
		"tag_name":    version.CurrentVersion.Original(),
		"ref":         baseBranch,
		"description": description,
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
		return fmt.Errorf("failed to create release. Status code: %d", resp.StatusCode)
	}

	fmt.Println("Release created successfully.")
	return nil
}

func (g Client) CheckRelease(version config.Versions) (bool, error) {
	url := fmt.Sprintf("%s/projects/%d/repository/tags", g.ApiURL, g.ProjectID)

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
		return false, fmt.Errorf("failed to fetch tags. Status code: %d", resp.StatusCode)
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
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
