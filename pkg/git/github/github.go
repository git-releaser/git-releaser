package github

import "github.com/thschue/git-releaser/pkg/config"

type Client struct {
	UserId      string
	AccessToken string
	ProjectURL  string
	Repository  string
	ApiURL      string
}

func (g Client) CreateRelease(baseBranch string, version config.Versions, description string) error {
	return nil
}

func (g Client) CheckRelease(version config.Versions) (bool, error) {
	return false, nil
}
