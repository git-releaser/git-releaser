package file

import (
	"github.com/Masterminds/semver"
	"github.com/git-releaser/git-releaser/pkg/config"
	"os"
	"testing"
)

const contentBetweenTagsInput = `<!-- x-git-releaser-version-start -->

Current version: 1.2.3

<!-- x-git-releaser-version-end -->`

const contentBetweenTagsOutput = `<!-- x-git-releaser-version-start -->

Current version: 2.0.0

<!-- x-git-releaser-version-end -->`

const contentWithAnnotationInput = `This is a Test with version 1.2.3 in the middle # x-git-releaser-version`

const contentWithAnnotationOutput = `This is a Test with version 2.0.0 in the middle # x-git-releaser-version`

func TestReplaceVersionBetweenTagsPositive(t *testing.T) {

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write initial content to the file
	if _, err := tempFile.Write([]byte(contentBetweenTagsInput)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Define the extra file config
	extraFile := config.ExtraFileConfig{
		Path: tempFile.Name(),
	}

	// Define the versions
	versions := config.Versions{
		NextVersion: *semver.MustParse("2.0.0"),
	}

	// Call the function
	if err := replaceVersionBetweenTags(extraFile, versions); err != nil {
		t.Fatalf("replaceVersionBetweenTags failed: %v", err)
	}

	// Read the modified content from the file
	modifiedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read from temp file: %v", err)
	}

	// Check the modified content
	if string(modifiedContent) != contentBetweenTagsOutput {
		t.Errorf("Unexpected content: got %v, want %v", string(modifiedContent), contentBetweenTagsOutput)
	}
}

func TestReplaceVersionLinesPositive(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write initial content to the file
	if _, err := tempFile.Write([]byte(contentWithAnnotationInput)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Define the extra file config
	extraFile := config.ExtraFileConfig{
		Path: tempFile.Name(),
	}

	// Define the versions
	versions := config.Versions{
		NextVersion: *semver.MustParse("2.0.0"),
	}

	// Call the function
	if err := replaceVersionLines(extraFile, versions); err != nil {
		t.Fatalf("replaceVersionLines failed: %v", err)
	}

	// Read the modified content from the file
	modifiedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read from temp file: %v", err)
	}

	// Check the modified content
	if string(modifiedContent) != contentWithAnnotationOutput {
		t.Errorf("Unexpected content: got %v, want %v", string(modifiedContent), contentWithAnnotationOutput)
	}
}
