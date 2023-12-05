package github

import "fmt"

type Client struct {
	UserId      string
	AccessToken string
	ProjectURL  string
}

func (g Client) CheckCreateBranch() (string, error) {
	fmt.Println("CheckCreateBranch")
	return "", fmt.Errorf("not implemented")
}

func (g Client) CheckCreatePullRequest(source string, target string) error {
	fmt.Println("CheckCreatePullRequest")
	return fmt.Errorf("not implemented")

}

func (g Client) CommitManifest(branchName string, content string) error {
	fmt.Println("CommitManifest")
	return fmt.Errorf("not implemented")
}
