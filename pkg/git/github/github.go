package github

type Client struct {
	UserId      string
	AccessToken string
	ProjectURL  string
	Repository  string
	ApiURL      string
}

func (g Client) CreateRelease(baseBranch string, version string, description string) error {
	return nil
}

func (g Client) CheckRelease(version string) (bool, error) {
	return false, nil
}
