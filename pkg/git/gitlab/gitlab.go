package gitlab

type Client struct {
	UserId      string
	AccessToken string
	ApiURL      string
	ProjectID   int
	ProjectURL  string
	DryRun      bool
}

func (g Client) GetHighestRelease() (string, error) {
	return "", nil
}
