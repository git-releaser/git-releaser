package changelog

import (
	"testing"
)

func TestParseConventionalCommits(t *testing.T) {
	commits := []Commit{
		{
			ID:        "abc123",
			Message:   "feat: Added new feature",
			Timestamp: "2022-01-01T00:00:00Z",
		},
		{
			ID:        "def456",
			Message:   "fix: Fixed a bug",
			Timestamp: "2022-01-02T00:00:00Z",
		},
		{
			ID:        "ghi789",
			Message:   "chore: Updated dependencies",
			Timestamp: "2022-01-03T00:00:00Z",
		},
		{
			ID:        "jkl123",
			Message:   "test123: Updated dependencies",
			Timestamp: "2022-01-03T00:00:00Z",
		},
	}

	expected := []ConventionalCommit{
		{
			Type:    "feat",
			Message: "Added new feature",
			ID:      "abc123",
		},
		{
			Type:    "fix",
			Message: "Fixed a bug",
			ID:      "def456",
		},
		{
			Type:    "chore",
			Message: "Updated dependencies",
			ID:      "ghi789",
		},
	}

	result := ParseCommits(commits)

	if len(result) != len(expected) {
		t.Errorf("Expected length: %d, got: %d", len(expected), len(result))
		return
	}

	for i, commit := range result {
		if commit.Type != expected[i].Type || commit.Message != expected[i].Message || commit.ID != expected[i].ID {
			t.Errorf("Expected commit: %+v, got: %+v", expected[i], commit)
		}
	}
}

func TestGenerateChangelog(t *testing.T) {
	commits := []ConventionalCommit{
		{
			Type:    "feat",
			Message: "Added new feature",
			ID:      "abc123",
		},
		{
			Type:    "fix",
			Message: "Fixed a bug",
			ID:      "def456",
		},
		{
			Type:    "chore",
			Message: "Updated dependencies",
			ID:      "ghi789",
		},
	}

	projectURL := "https://github.com/thschue/git-releaser"

	expected := `## Chores
- [Updated dependencies](https://github.com/thschue/git-releaser/commit/ghi789)

## Features
- [Added new feature](https://github.com/thschue/git-releaser/commit/abc123)

## Bug Fixes
- [Fixed a bug](https://github.com/thschue/git-releaser/commit/def456)

`

	result := GenerateChangelog(commits, projectURL)

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
