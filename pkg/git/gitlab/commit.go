package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/thschue/git-releaser/pkg/changelog"
	releaserconfig "github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/file"
	"net/http"
	"net/url"
)

// ConventionalCommit represents a conventional commit structure
type ConventionalCommit struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

func (g Client) CommitManifest(branchName string, content string, versions releaserconfig.Versions, extraFiles []releaserconfig.ExtraFileConfig) error {
	err := file.CommitManifest(branchName, g.UserId, g.AccessToken, content, versions, extraFiles)
	return err
}

func (g Client) GetCommitsSinceRelease(sinceRelease string) ([]changelog.Commit, error) {
	var giturl string
	var tagDate string
	var err error

	if sinceRelease != "0.0.0" && sinceRelease != "" {
		tagDate, err = g.getTagCommitDate(sinceRelease)
		if err != nil {
			fmt.Println("Could not get tag date: " + err.Error())
		}
	}

	if tagDate == "" {
		giturl = fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/repository/commits", g.ProjectID)
	} else {
		giturl = fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/repository/commits?since=%s", g.ProjectID, tagDate)
		fmt.Println(giturl)
	}

	req, err := http.NewRequest("GET", giturl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
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
	tagUrl := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/repository/tags/%s", g.ProjectID, url.PathEscape(tag))
	fmt.Println(tagUrl)
	tagReq, err := http.NewRequest("GET", tagUrl, nil)
	if err != nil {
		return "", err
	}

	tagReq.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	tagResp, err := client.Do(tagReq)
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
