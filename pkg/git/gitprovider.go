package git

import (
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/git/github"
	"github.com/git-releaser/git-releaser/pkg/git/gitlab"
	"log"
	"strconv"
	"strings"
)

type GitConfig struct {
	Provider           string
	UserId             string
	AccessToken        string
	ProjectUrl         string
	ApiUrl             string
	AdditionalConfig   map[string]string
	PropagationTargets []config.PropagationTarget
	DryRun             bool
}
type GitProvider interface {
	CheckCreateBranch(baseBranch string, targetVersion string, prefix string) (string, error)
	CheckCreateReleasePullRequest(source string, target string, versions config.Versions) error
	CommitManifest(branchName string, content string, versions config.Versions, extraFiles []config.ExtraFileConfig) error
	CreateRelease(baseBranch string, version config.Versions, description string) error
	CheckRelease(versions config.Versions) (bool, error)
	GetCommitsSinceRelease(version string) ([]changelog.Commit, error)
	GetHighestRelease() (string, error)
}

func NewGitClient(gitconfig GitConfig) GitProvider {
	if gitconfig.Provider == "" {
		gitconfig.Provider = "github"
	}

	switch strings.ToLower(gitconfig.Provider) {
	case "gitlab":
		if gitconfig.ApiUrl == "" {
			gitconfig.ApiUrl = "https://gitlab.com/api/v4"
		}

		projectID, err := strconv.Atoi(gitconfig.AdditionalConfig["projectId"])
		if err != nil {
			log.Fatal(err)
		}
		return &gitlab.Client{
			UserId:             gitconfig.UserId,
			AccessToken:        gitconfig.AccessToken,
			ApiURL:             gitconfig.ApiUrl,
			ProjectID:          projectID,
			ProjectURL:         gitconfig.ProjectUrl,
			PropagationTargets: gitconfig.PropagationTargets,
			DryRun:             gitconfig.DryRun,
		}

	case "github":
		if gitconfig.ApiUrl == "" {
			gitconfig.ApiUrl = "https://api.github.com"
		}

		fmt.Println(gitconfig.UserId)
		return github.NewClient(github.Client{
			UserId:             gitconfig.UserId,
			AccessToken:        gitconfig.AccessToken,
			ProjectURL:         gitconfig.ProjectUrl,
			Repository:         gitconfig.AdditionalConfig["repository"],
			ApiURL:             gitconfig.ApiUrl,
			PropagationTargets: gitconfig.PropagationTargets,
			DryRun:             gitconfig.DryRun,
		})
	}
	return nil
}
