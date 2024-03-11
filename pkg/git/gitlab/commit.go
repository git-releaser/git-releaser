package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	releaserconfig "github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/file"
	"net/http"
	"net/url"
)

func (g Client) CommitManifest(branchName string, content string, versions releaserconfig.Versions, extraFiles []releaserconfig.ExtraFileConfig) error {
	err := file.CommitManifest(branchName, g.UserId, g.AccessToken, content, versions, extraFiles, g.DryRun)
	return err
}

func (g Client) GetCommitsSinceRelease(sinceRelease string) ([]changelog.Commit, error) {
	var req GitLabRequest
	var tagDate string
	var err error

	if sinceRelease != "0.0.0" && sinceRelease != "" {
		tagDate, err = g.getTagCommitDate(sinceRelease)
		if err != nil {
			fmt.Println("Could not get tag date: " + err.Error())
		}
	}

	if tagDate == "" {
		req.URL = fmt.Sprintf("%s/projects/%d/repository/commits", g.ApiURL, g.ProjectID)
	} else {
		req.URL = fmt.Sprintf("%s/projects/%d/repository/commits?since=%s", g.ApiURL, g.ProjectID, tagDate)
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commits. Status code: %d", resp.StatusCode)
	}

	var commits []changelog.Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	return commits, nil
}

func (g Client) getTagCommitDate(tag string) (string, error) {
	req := GitLabRequest{
		URL: fmt.Sprintf("%s/projects/%d/repository/tags/%s", g.ApiURL, g.ProjectID, url.PathEscape(tag)),
	}
	tagResp, err := g.gitLabRequest(req)
	if err != nil {
		return "", err
	}
	defer tagResp.Body.Close()

	if tagResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get tag details. Status code: %d", tagResp.StatusCode)
	}

	var tagDetails struct {
		Commit struct {
			CommittedDate string `json:"created_at"`
		} `json:"commit"`
	}
	if err := json.NewDecoder(tagResp.Body).Decode(&tagDetails); err != nil {
		return "", err
	}
	return tagDetails.Commit.CommittedDate, nil
}
