package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"strings"
)

func GetGitHistory(path string, tag string) ([]object.Commit, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	// Get the tag reference
	ref, err := r.Tag(tag)
	if err != nil {
		return nil, err
	}

	// Resolve the tag to a commit
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// Get the commit history
	iter, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var commits []object.Commit
	// Iterate over the commit history
	err = iter.ForEach(func(c *object.Commit) error {
		// Only print the commits that are descendants of the tag
		if c.Committer.When.After(commit.Committer.When) && c.NumParents() > 1 {
			lines := strings.Split(c.Message, "\n")
			if len(lines) > 0 && strings.HasPrefix(lines[0], "Merge pull request") {
				lines = lines[1:]
			}
			c.Message = strings.Join(lines, " ")     // Join lines with a space instead of a newline
			c.Message = strings.TrimSpace(c.Message) // Remove leading and trailing white space
			commits = append(commits, *c)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, err
}
