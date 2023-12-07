package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/thschue/git-releaser/pkg/changelog"
	releaserconfig "github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/file"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// ConventionalCommit represents a conventional commit structure
type ConventionalCommit struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

func (g Client) CommitManifest(branchName string, content string, versions releaserconfig.Versions, extraFiles []releaserconfig.ExtraFileConfig) error {
	filePath := ".git-releaser-manifest.json"

	repository, err := git.PlainOpen("./")
	if err != nil {
		return err
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	// Create or update the file in the worktree
	err = os.WriteFile(filepath.Join(worktree.Filesystem.Root(), filePath), []byte(content), 0644)
	if err != nil {
		return err
	}

	// Add the file to the worktree
	_, err = worktree.Add(filePath)
	if err != nil {
		return err
	}

	for _, extraFile := range extraFiles {
		err = file.ReplaceVersionLines(extraFile, versions)
		if err != nil {
			fmt.Println("Could not update version in file: " + extraFile.Path)
		}

		err = file.ReplaceVersionBetweenTags(extraFile, versions)
		if err != nil {
			fmt.Println("Could not update version in file: " + extraFile.Path)
		}

		_, err = worktree.Add(extraFile.Path)
		if err != nil {
			fmt.Println("Could not add file to git: " + filepath.Join(worktree.Filesystem.Root(), extraFile.Path))
		}
	}

	// Commit the changes
	commit, err := worktree.Commit("releaser: update files for version "+versions.NextVersionSlug, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your.email@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// Update the branch reference to point to the new commit
	refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))
	err = repository.Storer.SetReference(plumbing.NewHashReference(refName, commit))
	if err != nil {
		return err
	}

	auth := &githttp.BasicAuth{
		Username: g.UserId,
		Password: g.AccessToken,
	}
	options := git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)),
		},
		Auth: auth,
	}

	// Push the changes to the remote repository
	err = repository.Push(&options)
	if err != nil {
		return err
	}

	return nil
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
