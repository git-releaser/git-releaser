package git

import (
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
		return &github.Client{
			UserId:      config.UserId,
			AccessToken: config.AccessToken,
			ProjectURL:  config.ProjectUrl,
		}
	}
	return nil
}
