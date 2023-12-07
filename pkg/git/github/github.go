package github

import (
	"context"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type Client struct {
	Context     context.Context
	UserId      string
	AccessToken string
	ProjectURL  string
	Repository  string
	ApiURL      string
	GHClient    *github.Client
}

func NewClient(client Client) Client {
	client.Context = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: client.AccessToken},
	)
	tc := oauth2.NewClient(client.Context, ts)

	client.GHClient = github.NewClient(tc)

	return client
}
