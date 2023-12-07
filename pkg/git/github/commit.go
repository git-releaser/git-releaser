package github

import (
	"encoding/json"
	"fmt"
	"github.com/thschue/git-releaser/pkg/changelog"
	releaserconfig "github.com/thschue/git-releaser/pkg/config"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

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

	// Commit the changes
	commit, err := worktree.Commit("chore: update version", &git.CommitOptions{
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
	var url string
	if sinceRelease == "0.1.0" || sinceRelease == "" {
		url = fmt.Sprintf("https://%s/repos/%s/commits", g.ApiURL, g.Repository)
	} else {
		url = fmt.Sprintf("%s/repos/%s/commits?since=%s", g.ApiURL, g.Repository, sinceRelease)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

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
