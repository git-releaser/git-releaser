package github

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	"github.com/thschue/git-releaser/pkg/changelog"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/naming"
	"sort"
	"strings"
)

func (g Client) CreateRelease(baseBranch string, version config.Versions, description string) error {
	highestRelease, err := g.GetHighestRelease()
	if err != nil {
		fmt.Println("github: could not get highest release")
	}
	commits, _ := g.GetCommitsSinceRelease(version.VersionPrefix + highestRelease)
	conventionalCommits := changelog.ParseConventionalCommits(commits)
	cl := changelog.GenerateChangelog(conventionalCommits, g.ProjectURL)

	if description == "" {
		description = naming.CreateReleaseDescription(version.CurrentVersionSlug, cl)
	}

	release := &github.RepositoryRelease{
		TagName:         github.String(version.CurrentVersionSlug),
		TargetCommitish: github.String(baseBranch),
		Name:            github.String("Release " + version.CurrentVersionSlug),
		Body:            github.String(description),
	}

	owner, repo := parseOwnerRepoFromURL(g.ProjectURL)
	_, _, err = g.GHClient.Repositories.CreateRelease(g.Context, owner, repo, release)
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

func (g Client) GetHighestRelease() (string, error) {
	var org string
	var repo string

	if len(strings.Split(g.Repository, "/")) == 2 {
		org = strings.Split(g.Repository, "/")[0]
		repo = strings.Split(g.Repository, "/")[1]
	}

	releases, _, err := g.GHClient.Repositories.ListReleases(g.Context, org, repo, nil)
	if err != nil {
		return "", err
	}

	if len(releases) == 0 {
		return "0.0.0", nil
	}

	versions := make([]*semver.Version, len(releases))
	for i, release := range releases {
		version, err := semver.NewVersion(release.GetTagName())
		if err != nil {
			continue // Ignore invalid versions
		}
		versions[i] = version
	}

	sort.Sort(semver.Collection(versions))

	return versions[len(versions)-1].String(), nil
}
