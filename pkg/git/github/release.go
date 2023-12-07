package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"github.com/thschue/git-releaser/pkg/config"
	"golang.org/x/oauth2"
	"strings"
)

func (g Client) CreateRelease(baseBranch string, version config.Versions, description string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	release := &github.RepositoryRelease{
		TagName:         github.String(version.CurrentVersionSlug),
		TargetCommitish: github.String(baseBranch),
		Name:            github.String("Release " + version.CurrentVersionSlug),
		Body:            github.String(description),
	}

	owner, repo := parseOwnerRepoFromURL(g.ProjectURL)
	_, _, err := client.Repositories.CreateRelease(ctx, owner, repo, release)
	if err != nil {
		return err
	}

	fmt.Println("Release created successfully.")
	return nil
}

func parseOwnerRepoFromURL(url string) (string, string) {
	// Assuming URL is of the form "https://github.com/owner/repo"
	parts := strings.Split(url, "/")
	return parts[len(parts)-2], parts[len(parts)-1]
}

func (g Client) CheckRelease(version config.Versions) (bool, error) {
	owner, repo := parseOwnerRepoFromURL(g.ProjectURL)
	tags, _, err := g.GHClient.Repositories.ListTags(g.Context, owner, repo, nil)
	if err != nil {
		return false, err
	}

	// Check if the desired tag is in the list
	for _, tag := range tags {
		if *tag.Name == version.CurrentVersionSlug {
			return true, nil
		}
	}

	return false, nil
}
