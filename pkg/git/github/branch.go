package github

import (
	"fmt"
	"github.com/google/go-github/v33/github"
	"strings"
)

func (g Client) CheckCreateBranch(baseBranch string, version string) (string, error) {
	branchName := fmt.Sprintf("release-%s", version)
	branchExists, _ := g.branchExists(branchName)
	if !branchExists {
		err := g.createBranch(baseBranch, branchName)
		if err != nil {
			return "", err
		}
	}
	return branchName, nil
}

func (g Client) branchExists(branchName string) (bool, error) {
	owner, repo := strings.Split(g.Repository, "/")[0], strings.Split(g.Repository, "/")[1]

	_, resp, err := g.GHClient.Repositories.GetBranch(g.Context, owner, repo, branchName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Branch does not exist
			return false, nil
		}
		// Other errors
		return false, err
	}

	// Branch exists
	return true, nil
}

// createBranch creates a branch in a GitHub repository
func (g Client) createBranch(baseBranch string, branchName string) error {
	owner, repo := strings.Split(g.Repository, "/")[0], strings.Split(g.Repository, "/")[1]

	// Get the SHA of the base branch
	baseRef, _, err := g.GHClient.Git.GetRef(g.Context, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		fmt.Println("Could not get the SHA of the base branch")
		return err
	}
	baseSha := baseRef.GetObject().GetSHA()

	// Create a new reference for the new branch
	ref := &github.Reference{
		Ref: github.String("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: github.String(baseSha),
		},
	}

	_, _, err = g.GHClient.Git.CreateRef(g.Context, owner, repo, ref)
	if err != nil {
		return err
	}

	fmt.Printf("Branch '%s' created successfully.\n", branchName)
	return nil
}
