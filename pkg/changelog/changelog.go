package changelog

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// ConventionalCommit represents a conventional commit structure
type ConventionalCommit struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func ParseConventionalCommits(commits []Commit) []ConventionalCommit {
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

// GenerateChangelog generates a changelog from conventional commits
func GenerateChangelog(commits []ConventionalCommit, projectURL string) string {
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
