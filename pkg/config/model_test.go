package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write initial content to the file
	initialContent := `target_branch: "main"
provider: "github"
extra_files: []
user_id: "test"
access_token: "test"
project_url: "https://github.com/test/test"
api_url: "https://api.github.com"
project_id: 1
versioning:
  version_prefix: "v"
  bump_minor_pre_major: true
  bump_patch_minor_pre_major: true`
	if _, err := tempFile.Write([]byte(initialContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Call the function
	config, err := ReadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	// Check the config
	if config.TargetBranch != "main" {
		t.Errorf("Unexpected TargetBranch: got %v, want %v", config.TargetBranch, "main")
	}
	if config.Provider != "github" {
		t.Errorf("Unexpected Provider: got %v, want %v", config.Provider, "github")
	}
}
