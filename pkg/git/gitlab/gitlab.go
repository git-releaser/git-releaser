package gitlab

type Client struct {
	UserId      string
	AccessToken string
	ApiURL      string
	ProjectID   int
	ProjectURL  string
}

func (g Client) GetHighestRelease() (string, error) {
	return "", nil
}
