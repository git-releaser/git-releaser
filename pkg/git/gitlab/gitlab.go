package gitlab

import "github.com/git-releaser/git-releaser/pkg/config"

type Client struct {
	UserId             string
	AccessToken        string
	ApiURL             string
	ProjectID          int
	ProjectURL         string
	PropagationTargets []config.PropagationTarget
	DryRun             bool
}

func (g Client) GetHighestRelease() (string, error) {
	return "", nil
}
