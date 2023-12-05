package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// ConventionalCommit represents a conventional commit structure
type ConventionalCommit struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

func (g Client) CommitManifest(branchName string, content string) error {
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

func (g Client) getCommitsSinceRelease(sinceRelease string) ([]Commit, error) {
	var url string
	if sinceRelease == "0.1.0" || sinceRelease == "" {
		url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/repository/commits", g.ProjectID)
	} else {
		url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/repository/commits?since=%s", g.ProjectID, sinceRelease)
	}
	req, err := http.NewRequest("GET", url, nil)
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

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	return commits, nil
}

func parseConventionalCommits(commits []Commit) []ConventionalCommit {
	var conventionalCommits []ConventionalCommit

	for _, commit := range commits {
		parts := strings.SplitN(commit.Message, ":", 2)
		if len(parts) == 2 {
			conventionalCommits = append(conventionalCommits, ConventionalCommit{
				Type:    strings.TrimSpace(parts[0]),
				Message: strings.TrimSpace(parts[1]),
				ID:      commit.ID,
			})
		}
	}

	return conventionalCommits
}

// generateChangelog generates a changelog from conventional commits
func generateChangelog(commits []ConventionalCommit, projectURL string) string {
	// Map to store commits grouped by type
	commitsByType := make(map[string][]ConventionalCommit)

	// Set to keep track of unique commit messages
	uniqueCommits := make(map[string]struct{})

	// Group commits by type and filter duplicates
	for _, commit := range commits {
		if _, exists := uniqueCommits[commit.Message]; !exists {
			commitsByType[commit.Type] = append(commitsByType[commit.Type], commit)
			uniqueCommits[commit.Message] = struct{}{}
		}
	}

	var changelogBuffer bytes.Buffer

	// Iterate over commit types in sorted order
	for _, commitType := range getSortedKeys(commitsByType) {
		// Add heading for the commit type
		changelogBuffer.WriteString(fmt.Sprintf("## %s\n", commitType))

		// Iterate over commits for the current type
		for _, commit := range commitsByType[commitType] {
			// Add a link to the commit
			commitLink := fmt.Sprintf("[%s](%s/commit/%s)", commit.Message, projectURL, commit.ID)
			changelogBuffer.WriteString(fmt.Sprintf("- %s\n", commitLink))
		}

		changelogBuffer.WriteString("\n") // Add a newline between sections
	}

	return changelogBuffer.String()
}

func getSortedKeys(m map[string][]ConventionalCommit) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
