package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Release struct {
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

func (g Client) CreateRelease(version string, description string) error {
	url := fmt.Sprintf("%s/projects/%d/releases", g.ApiURL, g.ProjectID)

	payload := map[string]interface{}{
		"tag_name":    version,
		"ref":         "main",
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
