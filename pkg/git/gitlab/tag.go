package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Tag struct {
	Name   string `json:"name"`
	Commit struct {
		ID string `json:"id"`
	} `json:"commit"`
}

func (g Client) CheckRelease(version string) (bool, error) {
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
		if tag.Name == version {
			return true, nil
		}
	}
	return false, nil
}
