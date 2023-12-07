package file

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/thschue/git-releaser/pkg/config"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func CommitManifest(branchName string, userid string, token string, content string, versions config.Versions, extraFiles []config.ExtraFileConfig) error {
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
		err = replaceVersionLines(extraFile, versions)
		if err != nil {
			fmt.Println("Could not update version in file: " + extraFile.Path)
		}

		err = replaceVersionBetweenTags(extraFile, versions)
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
			Name:  "git-releaser",
			Email: "no-reply@git-releaser.com",
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
		Username: userid,
		Password: token,
	}
	options := git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []gitconfig.RefSpec{
			gitconfig.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)),
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

func replaceVersionLines(extraFile config.ExtraFileConfig, versions config.Versions) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		return err
	}

	// Define a regular expression to match the version string with the annotation format
	versionRegex := regexp.MustCompile(`(?m)(.*?)(\d+\.\d+\.\d+)(.*?)# x-git-releaser-version`)

	// Replace all occurrences of the version in annotated lines with the new version
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}"+versions.NextVersion.String()+"${3}# x-git-releaser-version")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func replaceVersionBetweenTags(extraFile config.ExtraFileConfig, versions config.Versions) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		return err
	}

	// Define a regular expression to match the version string and the rest of the text between the start and end tags
	versionRegex := regexp.MustCompile(`(?s)(<!-- x-git-releaser-version-start -->)(.*?)(\d+\.\d+\.\d+)(.*?)(<!-- x-git-releaser-version-end -->)`)

	// Replace the version string with the new version, preserving the rest of the text
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}${2}"+versions.NextVersion.String()+"${4}${5}")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
