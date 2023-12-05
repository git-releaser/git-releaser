package gitlab

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"path/filepath"
	"time"
)

func (g GitLabClient) CommitManifest(branchName string, content string) error {
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

	fmt.Println(filepath.Join(worktree.Filesystem.Root(), filePath))

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

	auth := &http.BasicAuth{
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
