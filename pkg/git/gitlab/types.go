package gitlab

type MergeRequest struct {
	ID           int      `json:"id"`
	IID          int      `json:"iid"`
	SourceBranch string   `json:"source_branch"`
	TargetBranch string   `json:"target_branch"`
	Title        string   `json:"title"`
	Labels       []string `json:"labels"`
	State        string   `json:"state"`
}
