package gitlab

import (
	"bytes"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/git/common"
	"io"
	"net/http"
)

type Client struct {
	UserId             string
	AccessToken        string
	ApiURL             string
	ProjectID          int
	ProjectURL         string
	PropagationTargets []config.PropagationTarget
	ConfigUpdates      []config.ConfigUpdate
	DryRun             bool
	GoGitConfig        common.GoGitRepository
}

type Request struct {
	Method  string
	URL     string
	Payload []byte
}

type Response struct {
	StatusCode int
	Body       []byte
}

func (g Client) GetHighestRelease() (string, error) {
	return "", nil
}

func (g Client) gitLabRequest(request Request) (Response, error) {
	var req *http.Request
	var err error

	if request.Method == "" {
		request.Method = "GET"
	}

	switch request.Method {
	case "PUT", "POST":
		req, err = http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Payload))
		if err != nil {
			return Response{}, err
		}
	case "GET", "DELETE":
		req, err = http.NewRequest(request.Method, request.URL, nil)
		if err != nil {
			return Response{}, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	// Read the body content into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{StatusCode: resp.StatusCode}, err
	}

	defer resp.Body.Close()

	return Response{StatusCode: resp.StatusCode, Body: bodyBytes}, err
}
