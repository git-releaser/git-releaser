package config

type ConventionalCommit struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Message string `json:"message"`
	ID      string `json:"id"`
}
