package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/changelog"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/helpers"
	"github.com/git-releaser/git-releaser/pkg/naming"
	"net/http"
)

func (g Client) CheckCreateReleasePullRequest(source string, target string, versions config.Versions) error {
	// Check if a pull request with the same source and target branches already exists
	existingPR, err := g.getMergeRequestBySourceAndTarget(source, target)
	if err != nil {
		return err
	}

	commits, _ := g.GetCommitsSinceRelease(versions.CurrentVersion.Original())
	conventionalCommits := changelog.ParseCommits(commits)
	cl := changelog.GenerateChangelog(conventionalCommits, g.ProjectURL)

	m := MergeRequest{
		SourceBranch: source,
		TargetBranch: target,
		Title:        naming.GeneratePrTitle(versions.NextVersion.Original()),
		Description:  naming.CreatePrDescription(versions.NextVersion.Original(), cl, g.PropagationTargets, g.ConfigUpdates),
		Labels:       []string{"release"},
	}

	if existingPR.IID != 0 {
		// If the pull request already exists, update its description
		err := existingPR.Update(g)
		if err != nil {
			return err
		}

	} else {
		err := m.Create(g)
		if err != nil {
			return err
		}
	}

	// Check if other git-releaser pull requests exist and close them
	err = g.closeOldPullRequests(source)
	if err != nil {
		return err
	}

	return nil
}

func (g Client) CheckCreateFileMergeRequest(source string, target string) error {
	fmt.Println("Checking if a pull request for the file update already exists")
	// Check if a pull request with the same source and target branches already exists
	existingPR, err := g.getMergeRequestBySourceAndTarget(source, target)
	if err != nil {
		return err
	}

	m := MergeRequest{
		SourceBranch: source,
		TargetBranch: target,
		Title:        fmt.Sprintf("Updating %s to %s", source, target),
		Labels:       []string{"release-updates"},
	}

	fmt.Println("Checking if pull request already exists")
	if existingPR.IID != 0 {
		// If the pull request already exists, update its description
		fmt.Println("Pull request already exists, will update it")
		err := existingPR.Update(g)
		if err != nil {
			return err
		}

	} else {
		fmt.Println("Pull request does not exist, will create it")
		err := m.Create(g)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g Client) closeOldPullRequests(currentSource string) error {
	mergeRequests, err := g.getMergeRequests()
	if err != nil {
		return err
	}

	for _, mr := range mergeRequests {
		// Check if the merge request is open and has a "release" label
		if mr.State == "opened" && helpers.Contains(mr.Labels, "release") && mr.SourceBranch != currentSource {
			// Close the merge request
			err := mr.Close(g)
			if err != nil {
				return err
			}

			// Delete the source branch
			err = g.deleteBranch(mr.SourceBranch)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m MergeRequest) Close(g Client) error {
	var err error
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/merge_requests/%d", g.ApiURL, g.ProjectID, m.IID),
		Method: "PUT",
	}

	payload := map[string]interface{}{
		"state_event": "close",
	}

	req.Payload, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to close merge request. Status code: %d, Body: %s", resp.StatusCode, resp.Body)
	}

	return nil
}

func (m MergeRequest) Create(g Client) error {
	var err error

	req := Request{
		URL: fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID),
	}

	payload := map[string]interface{}{
		"source_branch": m.SourceBranch,
		"target_branch": m.TargetBranch,
		"title":         m.Title,
		"description":   m.Description,
		"labels":        m.Labels,
	}

	req.Payload, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	req.Method = "POST"

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create/update pull request. Status code: %d, Body: %s", resp.StatusCode, resp.Body)
	}

	fmt.Println("Pull request created successfully.")

	return nil
}

func (m MergeRequest) Update(g Client) error {
	var err error

	req := Request{
		URL: fmt.Sprintf("%s/projects/%d/merge_requests/%d", g.ApiURL, g.ProjectID, m.IID),
	}

	payload := map[string]interface{}{
		"source_branch": m.SourceBranch,
		"target_branch": m.TargetBranch,
		"title":         m.Title,
		"description":   m.Description,
		"labels":        m.Labels,
	}

	req.Payload, err = json.Marshal(payload)
	if err != nil {
		return err
	}

	req.Method = "PUT"

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create/update pull request. Status code: %d, Body: %s", resp.StatusCode, resp.Body)
	}

	fmt.Println("Pull request updated successfully.")

	return nil
}

func (g Client) getMergeRequestBySourceAndTarget(source, target string) (MergeRequest, error) {
	mergeRequests, err := g.getMergeRequests()
	if err != nil {
		return MergeRequest{IID: 0}, err
	}

	// Find the ID of the existing pull request with the same source and target branches
	for _, pr := range mergeRequests {
		if pr.SourceBranch == source && pr.TargetBranch == target && pr.State == "opened" {
			if pr.IID != 0 {
				return pr, nil
			}
			return pr, nil
		}
	}

	return MergeRequest{IID: 0}, nil // No existing pull request found
}

func (g Client) getMergeRequests() ([]MergeRequest, error) {
	req := Request{
		URL:    fmt.Sprintf("%s/projects/%d/merge_requests", g.ApiURL, g.ProjectID),
		Method: "GET",
	}

	resp, err := g.gitLabRequest(req)
	if err != nil {
		return []MergeRequest{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []MergeRequest{}, fmt.Errorf("failed to fetch merge requests. Status code: %d, Body: %s", resp.StatusCode, resp.Body)
	}

	var mergeRequests []MergeRequest

	if err := json.Unmarshal(resp.Body, &mergeRequests); err != nil {
		return []MergeRequest{}, err
	}

	return mergeRequests, nil // No existing pull request found
}
