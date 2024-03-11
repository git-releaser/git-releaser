package gitlab

import (
	"bytes"
	"github.com/git-releaser/git-releaser/pkg/config"
	"net/http"
)

type Client struct {
	UserId             string
	AccessToken        string
	ApiURL             string
	ProjectID          int
	ProjectURL         string
	PropagationTargets []config.PropagationTarget
	DryRun             bool
}

type GitLabRequest struct {
	Method  string
	URL     string
	Payload []byte
}

func (g Client) GetHighestRelease() (string, error) {
	return "", nil
}

func (g Client) gitLabRequest(request GitLabRequest) (*http.Response, error) {
	var req *http.Request
	var err error

	if request.Method == "" {
		request.Method = "GET"
	}

	switch request.Method {
	case "PUT", "POST":
		req, err = http.NewRequest("PUT", request.URL, bytes.NewBuffer(request.Payload))
		if err != nil {
			return &http.Response{}, err
		}
	case "GET":
		req, err = http.NewRequest("GET", request.URL, nil)
		if err != nil {
			return &http.Response{}, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, err
}
