package gitlab

type GitLabClient struct {
	UserId      string
	AccessToken string
	ApiURL      string
	ProjectID   int
	ProjectURL  string
}
