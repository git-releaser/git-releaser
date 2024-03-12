package github

import (
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	releaserconfig "github.com/git-releaser/git-releaser/pkg/config"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
)

func (g Client) CommitManifest(branchName string, content string, versions releaserconfig.Versions, extraFiles []releaserconfig.ExtraFileConfig) error {
	err := g.GoGitConfig.CommitManifest(branchName, content, versions, extraFiles, g.DryRun)
	return err
}

func (g Client) CommitFile(branchName string, content string, fileName string) error {
	err := g.GoGitConfig.CommitFile(branchName, content, fileName)
	return err
}

func (g Client) GetCommitsSinceRelease(sinceRelease string) ([]changelog.Commit, error) {
	var org string
	var repo string

	opt := &github.CommitsListOptions{}

	if len(strings.Split(g.Repository, "/")) == 2 {
		org = strings.Split(g.Repository, "/")[0]
		repo = strings.Split(g.Repository, "/")[1]
	}

	_, tagDate, err := g.getTagCommitSHA(org, repo, sinceRelease)
	if err != nil {
		fmt.Println("github: could not get tag commit SHA")
	}

	if tagDate != nil {
		opt.Since = *tagDate
	}

	ghCommits, _, err := g.GHClient.Repositories.ListCommits(g.Context, org, repo, opt)
	if err != nil {
		return nil, err
	}

	var commits []changelog.Commit
	for _, ghCommit := range ghCommits {
		date := ghCommit.Commit.Author.Date
		commit := changelog.Commit{
			ID:        *ghCommit.SHA,
			Message:   *ghCommit.Commit.Message,
			Timestamp: date.String(),
		}
		commits = append(commits, commit)
	}
	return commits, nil
}

func (g Client) getTagCommitSHA(owner string, repo string, tagName string) (string, *time.Time, error) {
	tags, _, err := g.GHClient.Repositories.ListTags(g.Context, owner, repo, nil)
	if err != nil {
		return "", nil, err
	}

	for _, tag := range tags {
		if *tag.Name == tagName {
			commit, _, err := g.GHClient.Repositories.GetCommit(g.Context, owner, repo, *tag.Commit.SHA)
			if err != nil {
				return "", nil, err
			}
			return *tag.Commit.SHA, commit.Commit.Author.Date, nil
		}
	}

	return "", nil, fmt.Errorf("tag not found")
}
