package changelog

import (
	"bytes"
	"fmt"
	"slices"
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

var validCommitTypes = map[string]string{
	"feat":  "Features",
	"fix":   "Bug Fixes",
	"chore": "Chores",
	"docs":  "Documentation",
	// Add more types as needed
}

func ParseCommits(commits []Commit) []ConventionalCommit {
	var conventionalCommits []ConventionalCommit
	var commitTypes []string
	for k := range validCommitTypes {
		commitTypes = append(commitTypes, k)
	}

	for _, commit := range commits {
		parts := strings.SplitN(commit.Message, ":", 2)
		if len(parts) == 2 {
			message := strings.Split(strings.TrimSpace(parts[1]), "\n")[0]
			if slices.Contains(commitTypes, strings.TrimSpace(parts[0])) {
				message := strings.Split(strings.TrimSpace(parts[1]), "\n")[0]
				conventionalCommits = append(conventionalCommits, ConventionalCommit{
					Type:    strings.TrimSpace(parts[0]),
					Message: message,
					ID:      commit.ID,
				})
			} else {
				conventionalCommits = append(conventionalCommits, ConventionalCommit{
					Type:    "other",
					Message: message,
					ID:      commit.ID,
				})
			}
		}
	}
	return conventionalCommits
}

// GenerateChangelog generates a changelog from conventional commits
func GenerateChangelog(commits []ConventionalCommit, projectURL string) string {
	// Map to store commits grouped by type
	commitsByType := make(map[string][]ConventionalCommit)
	// Slice to store other types of commits
	var otherCommits []ConventionalCommit

	// Set to keep track of unique commit messages
	uniqueCommits := make(map[string]struct{})

	// Group commits by type and filter duplicates
	for _, commit := range commits {
		if _, exists := uniqueCommits[commit.Message]; !exists {
			if _, exists := validCommitTypes[commit.Type]; exists {
				commitsByType[commit.Type] = append(commitsByType[commit.Type], commit)
			} else {
				otherCommits = append(otherCommits, commit)
			}
			uniqueCommits[commit.Message] = struct{}{}
		}
	}

	var changelogBuffer bytes.Buffer

	// Iterate over commit types in sorted order
	for _, commitType := range getSortedKeys(commitsByType) {
		// Add heading for the commit type
		changelogBuffer.WriteString(fmt.Sprintf("## %s\n", validCommitTypes[commitType]))

		// Iterate over commits for the current type
		for _, commit := range commitsByType[commitType] {
			// Add a link to the commit
			commitLink := fmt.Sprintf("[%s](%s/commit/%s)", commit.Message, projectURL, commit.ID)
			changelogBuffer.WriteString(fmt.Sprintf("- %s\n", commitLink))
		}

		changelogBuffer.WriteString("\n") // Add a newline between sections
	}

	// Add "Others" section if there are any other commits
	if len(otherCommits) > 0 {
		changelogBuffer.WriteString("## Others\n")
		for _, commit := range otherCommits {
			// Add a link to the commit
			commitLink := fmt.Sprintf("[%s](%s/commit/%s)", commit.Message, projectURL, commit.ID)
			changelogBuffer.WriteString(fmt.Sprintf("- %s\n", commitLink))
		}
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
