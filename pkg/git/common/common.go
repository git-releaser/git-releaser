package common

import (
	"errors"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type GoGitRepository struct {
	RepositoryUrl string
	Auth          *githttp.BasicAuth
	Repository    *git.Repository
	Worktree      *git.Worktree
}

type ChangeSet struct {
	fileName string
	content  string
}

func (g *GoGitRepository) CheckoutBranch(target string) error {
	var err error

	g.Repository = &git.Repository{}
	switch target {
	case "plain":
		fs := osfs.New(".")
		g.Repository, err = git.PlainOpen(fs.Root())
		if err != nil {
			return err
		}
	case "temp":
		storer := memory.NewStorage()
		fs := memfs.New()

		g.Repository, err = git.Clone(storer, fs, &git.CloneOptions{
			URL:  g.RepositoryUrl,
			Auth: g.Auth,
		})
		if err != nil {
			return err
		}
	}

	g.Worktree, err = g.Repository.Worktree()
	if err != nil {
		return err
	}

	err = g.Worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       g.Auth,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		fmt.Println("Could not pull the latest changes")
		return err
	}
	return nil
}

func (g GoGitRepository) CommitFile(branchName string, changeset []ChangeSet) error {
	if g.Worktree == nil {
		err := g.CheckoutBranch("temp")
		if err != nil {
			fmt.Println("Could not checkout branch")
			return err
		}
	}

	fmt.Println(g.Worktree)

	_ = g.Worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName)),
		Create: true,
	})

	for _, change := range changeset {

		file, err := g.Worktree.Filesystem.Create(change.fileName)
		if err != nil {
			fmt.Println("Could not create file: "+g.Worktree.Filesystem.Root(), change.fileName)
			continue
		}
		defer file.Close()

		_, err = file.Write([]byte(change.content))
		if err != nil {
			fmt.Println("Could not write to file: "+g.Worktree.Filesystem.Root(), change.fileName)
			continue
		}

		// Add the file to the worktree
		_, err = g.Worktree.Add(change.fileName)
		if err != nil {
			fmt.Println("Could not add file to git: " + filepath.Join(g.Worktree.Filesystem.Root(), change.fileName))
			continue
		}
	}

	// Commit the changes
	commit, err := g.Worktree.Commit("releaser: update files", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "git-releaser",
			Email: "no-reply@git-releaser.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Println("Could not commit changes")
		return err
	}

	// Update the branch reference to point to the new commit
	refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))
	err = g.Repository.Storer.SetReference(plumbing.NewHashReference(refName, commit))
	if err != nil {
		return err
	}

	options := git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []gitconfig.RefSpec{
			gitconfig.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)),
		},
		Auth:  g.Auth,
		Force: true,
	}

	// Push the changes to the remote repository
	err = g.Repository.Push(&options)
	if err != nil {
		fmt.Println("Could not push the changes")
		return err
	}
	return nil
}

func (g GoGitRepository) CommitManifest(branchName string, content string, versions config.Versions, extraFiles []config.ExtraFileConfig, dryRun bool) error {
	if g.Worktree == nil {
		err := g.CheckoutBranch("plain")
		if err != nil {
			return err
		}
	}

	filePath := ".git-releaser-manifest.json"

	// Create or update the file in the worktree

	file, err := g.Worktree.Filesystem.Create(filePath)
	if err != nil {
		fmt.Println("Could not create file: "+g.Worktree.Filesystem.Root(), filePath)
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Println("Could not write to file: "+g.Worktree.Filesystem.Root(), filePath)
		return err
	}

	// Add the file to the worktree
	_, err = g.Worktree.Add(filePath)
	if err != nil {
		fmt.Println("Could not add file to git: " + filepath.Join(g.Worktree.Filesystem.Root(), filePath))
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

		_, err = g.Worktree.Add(extraFile.Path)
		if err != nil {
			fmt.Println("Could not add file to git: " + filepath.Join(g.Worktree.Filesystem.Root(), extraFile.Path))
		}
	}

	if dryRun {
		fmt.Println("Dry run: would commit and push changes")
		return nil
	}
	// Commit the changes
	commit, err := g.Worktree.Commit("releaser: update files for version "+versions.NextVersion.Original(), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "git-releaser",
			Email: "no-reply@git-releaser.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Println("Could not commit changes")
		return err
	}

	// Update the branch reference to point to the new commit
	refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))
	err = g.Repository.Storer.SetReference(plumbing.NewHashReference(refName, commit))
	if err != nil {
		return err
	}

	options := git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []gitconfig.RefSpec{
			gitconfig.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)),
		},
		Auth:  g.Auth,
		Force: true,
	}

	// Push the changes to the remote repository
	err = g.Repository.Push(&options)
	if err != nil {
		fmt.Println("Could not push the changes")
		return err
	}

	return nil
}

func replaceVersionLines(extraFile config.ExtraFileConfig, versions config.Versions) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		fmt.Println("Could not read file: " + extraFile.Path)
		return err
	}

	// Define a regular expression to match the version string with the annotation format
	versionRegex := regexp.MustCompile(`(?m)(.*?)(\d+\.\d+\.\d+)(.*?)# x-git-releaser-version`)

	// Replace all occurrences of the version in annotated lines with the new version
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}"+versions.NextVersion.String()+"${3}# x-git-releaser-version")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		fmt.Println("Could not write file: " + extraFile.Path)
		return err
	}

	return nil
}

func replaceVersionBetweenTags(extraFile config.ExtraFileConfig, versions config.Versions) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		fmt.Println("Could not read file: " + extraFile.Path)
		return err
	}

	// Define a regular expression to match the version string and the rest of the text between the start and end tags
	versionRegex := regexp.MustCompile(`(?s)(<!-- x-git-releaser-version-start -->)(.*?)(\d+\.\d+\.\d+)(.*?)(<!-- x-git-releaser-version-end -->)`)

	// Replace the version string with the new version, preserving the rest of the text
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}${2}"+versions.NextVersion.String()+"${4}${5}")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		fmt.Println("Could not write file: " + extraFile.Path)
		return err
	}

	return nil
}

func (g GoGitRepository) ReplaceTaggedLines(filenames []string, sourceTag string, replaceTag string) ([]ChangeSet, error) {
	var changes []ChangeSet

	if g.Worktree == nil {
		err := g.CheckoutBranch("temp")
		if err != nil {
			return []ChangeSet{}, err
		}
	}

	for _, filename := range filenames {
		file, err := g.Worktree.Filesystem.Open(filename)
		if err != nil {
			fmt.Println("Could not open file: " + filename)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Could not read file: " + filename)
		}

		// Define a regular expression to match the version string with the annotation format
		versionRegex := regexp.MustCompile(`(?m)(.*?)(\d+\.\d+\.\d+)(.*?)# x-git-releaser:` + sourceTag)

		// Replace all occurrences of the version in annotated lines with the new version
		modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}"+replaceTag+"${3}# x-git-releaser:"+sourceTag)

		changes = append(changes, ChangeSet{fileName: filename, content: modifiedContent})
	}

	return changes, nil
}
