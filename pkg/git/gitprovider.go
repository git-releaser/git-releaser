package git

import (
	"fmt"
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
	AdditionalConfig map[string]string
}
type GitProvider interface {
	CheckCreateBranch(targetVersion string) (string, error)
	CheckCreatePullRequest(source string, target string, currentVersion string, targetVersion string) error
	CommitManifest(branchName string, content string) error
}

func NewGitClient(config GitConfig) GitProvider {
	if config.Provider == "" {
		config.Provider = "github"
	}

	switch strings.ToLower(config.Provider) {
	case "gitlab":
		projectID, err := strconv.Atoi(config.AdditionalConfig["projectId"])
		if err != nil {
			log.Fatal(err)
		}
		return &gitlab.Client{
			UserId:      config.UserId,
			AccessToken: config.AccessToken,
			ApiURL:      config.AdditionalConfig["apiUrl"],
			ProjectID:   projectID,
			ProjectURL:  config.ProjectUrl,
		}

	case "github":
		if config.AdditionalConfig["apiUrl"] == "" {
			config.AdditionalConfig["apiUrl"] = "https://api.github.com"
		}

		fmt.Println(config.UserId)
		return &github.Client{
			UserId:      config.UserId,
			AccessToken: config.AccessToken,
			ProjectURL:  config.ProjectUrl,
			Repository:  config.AdditionalConfig["repository"],
			ApiURL:      config.AdditionalConfig["apiUrl"],
		}
	}
	return nil
}
