package github

import (
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/naming"
	"github.com/google/go-github/v33/github"
	"strings"
)

func (g Client) CheckCreatePullRequest(source string, target string, versions config.Versions) error {
	err := g.createPullRequest(source, target, versions)
	if err != nil {
		return err
	}
	return nil
}

func (g Client) createPullRequest(source string, target string, versions config.Versions) error {
	owner, repo := strings.Split(g.Repository, "/")[0], strings.Split(g.Repository, "/")[1]

	// Check if a pull request with the same source and target branches already exists
	existingPrNumber, err := g.getExistingPullRequestNumber(source, target)
	if err != nil {
		return err
	}

	commits, _ := g.GetCommitsSinceRelease(versions.CurrentVersion.Original())
	conventionalCommits := changelog.ParseConventionalCommits(commits)
	cl := changelog.GenerateChangelog(conventionalCommits, g.ProjectURL)

	title := naming.GeneratePrTitle(versions.NextVersion.Original())
	description := naming.CreatePrDescription(versions.NextVersion.Original(), cl)

	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(description),
		Head:  github.String(source),
		Base:  github.String(target),
	}

	if existingPrNumber != 0 {
		// If the pull request already exists, update its description
		existingPr, _, err := g.GHClient.PullRequests.Get(g.Context, owner, repo, existingPrNumber)
		if err != nil {
			return err
		}
		existingPr.Title = newPR.Title
		existingPr.Body = newPR.Body

		if g.DryRun {
			fmt.Println("Pull request already exists, would update it")
			return nil
		}
		_, _, err = g.GHClient.PullRequests.Edit(g.Context, owner, repo, existingPrNumber, existingPr)
		if err != nil {
			return err
		}
		fmt.Println("Pull request updated successfully.")
	} else {
		// If the pull request doesn't exist, create a new one
		if g.DryRun {
			fmt.Println("Dry run: pull request would be created with the following details:")
			fmt.Println("Title: " + title)
			fmt.Println("Description: " + description)
			fmt.Println("Source branch: " + source)
			fmt.Println("Target branch: " + target)
			return nil
		}

		_, response, err := g.GHClient.PullRequests.Create(g.Context, owner, repo, newPR)
		if err != nil {
			if response.StatusCode == 403 {
				fmt.Println("Could not create pull request: " + err.Error())
				fmt.Println("Please make sure that the access token has the 'repo' scope.")
			}
			return err
		}
		fmt.Println("Pull request created successfully.")
	}

	return nil
}

func (g Client) getExistingPullRequestNumber(source, target string) (int, error) {
	owner, repo := strings.Split(g.Repository, "/")[0], strings.Split(g.Repository, "/")[1]

	// Fetch all pull requests
	opts := &github.PullRequestListOptions{
		State: "open",
	}
	pullRequests, _, err := g.GHClient.PullRequests.List(g.Context, owner, repo, opts)
	if err != nil {
		return 0, err
	}

	// Find the number of the existing pull request with the same source and target branches
	for _, pr := range pullRequests {
		if pr.GetHead().GetRef() == source && pr.GetBase().GetRef() == target {
			return pr.GetNumber(), nil
		}
	}

	return 0, nil // No existing pull request found
}
