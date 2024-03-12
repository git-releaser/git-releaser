package git

import (
	"github.com/Masterminds/semver"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/git/common"
	"github.com/git-releaser/git-releaser/pkg/git/github"
	"github.com/git-releaser/git-releaser/pkg/git/gitlab"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	Provider           string
	UserId             string
	AccessToken        string
	ProjectUrl         string
	ApiUrl             string
	AdditionalConfig   map[string]string
	PropagationTargets []config.PropagationTarget
	ConfigUpdates      []config.ConfigUpdate
	DryRun             bool
}
type Provider interface {
	CheckCreateBranch(baseBranch string, targetVersion string, prefix string) (string, error)
	CheckCreateReleasePullRequest(source string, target string, versions config.Versions) error
	CheckCreateFileMergeRequest(source string, target string) error
	CommitManifest(branchName string, content string, versions config.Versions, extraFiles []config.ExtraFileConfig) error
	CommitFile(branchName string, content string, fileName string) error
	CreateRelease(baseBranch string, version config.Versions, description string) error
	CheckRelease(versions config.Versions) (bool, error)
	GetCommitsSinceRelease(version string) ([]changelog.Commit, error)
	GetHighestRelease() (semver.Version, error)
}

func NewGitClient(gitconfig Config) Provider {
	if gitconfig.Provider == "" {
		gitconfig.Provider = "github"
	}

	goGitConfig := common.GoGitRepository{
		RepositoryUrl: gitconfig.ProjectUrl,
		Auth: &githttp.BasicAuth{
			Username: gitconfig.UserId,
			Password: gitconfig.AccessToken,
		},
	}

	switch strings.ToLower(gitconfig.Provider) {
	case "gitlab":
		if gitconfig.ApiUrl == "" {
			gitconfig.ApiUrl = "https://gitlab.com/api/v4"
		}

		projectID, err := strconv.Atoi(gitconfig.AdditionalConfig["projectId"])
		if err != nil {
			log.Fatal("Could not convert ProjectID to int. Please check your configuration file.")
		}

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
			GoGitConfig:        goGitConfig,
			ConfigUpdates:      gitconfig.ConfigUpdates,
			DryRun:             gitconfig.DryRun,
		}

	case "github":
		if gitconfig.ApiUrl == "" {
			gitconfig.ApiUrl = "https://api.github.com"
		}

		return github.NewClient(github.Client{
			UserId:             gitconfig.UserId,
			AccessToken:        gitconfig.AccessToken,
			ProjectURL:         gitconfig.ProjectUrl,
			Repository:         gitconfig.AdditionalConfig["repository"],
			ApiURL:             gitconfig.ApiUrl,
			PropagationTargets: gitconfig.PropagationTargets,
			GoGitConfig:        goGitConfig,
			DryRun:             gitconfig.DryRun,
		})
	}
	return nil
}
