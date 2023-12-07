package git

import (
	"fmt"
	"github.com/thschue/git-releaser/pkg/changelog"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/git/github"
	"github.com/thschue/git-releaser/pkg/git/gitlab"
	"log"
	"strconv"
	"strings"
)

type GitConfig struct {
	Provider         string
	UserId           string
	AccessToken      string
	ProjectUrl       string
	ApiUrl           string
	AdditionalConfig map[string]string
}
type GitProvider interface {
	CheckCreateBranch(baseBranch string, targetVersion string) (string, error)
	CheckCreatePullRequest(source string, target string, versions config.Versions) error
	CommitManifest(branchName string, content string, versions config.Versions, extraFiles []config.ExtraFileConfig) error
	CreateRelease(baseBranch string, version config.Versions, description string) error
	CheckRelease(versions config.Versions) (bool, error)

	GetCommitsSinceRelease(version string) ([]changelog.Commit, error)
}

func NewGitClient(config GitConfig) GitProvider {
	if config.Provider == "" {
		config.Provider = "github"
	}

	switch strings.ToLower(config.Provider) {
	case "gitlab":
		if config.ApiUrl == "" {
			config.ApiUrl = "https://gitlab.com/api/v4"
		}

		projectID, err := strconv.Atoi(config.AdditionalConfig["projectId"])
		if err != nil {
			log.Fatal(err)
		}
		return &gitlab.Client{
			UserId:      config.UserId,
			AccessToken: config.AccessToken,
			ApiURL:      config.ApiUrl,
			ProjectID:   projectID,
			ProjectURL:  config.ProjectUrl,
		}

	case "github":
		if config.ApiUrl == "" {
			config.ApiUrl = "https://api.github.com"
		}

		fmt.Println(config.UserId)
		return &github.Client{
			UserId:      config.UserId,
			AccessToken: config.AccessToken,
			ProjectURL:  config.ProjectUrl,
			Repository:  config.AdditionalConfig["repository"],
			ApiURL:      config.ApiUrl,
		}
	}
	return nil
}
