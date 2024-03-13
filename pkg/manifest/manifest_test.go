package manifest

import (
	"github.com/Masterminds/semver"
	"github.com/git-releaser/git-releaser/pkg/naming"
	"os"
	"testing"
)

func TestGetCurrentVersion(t *testing.T) {
	// Create a temporary file as a mock manifest file
	tempFile, err := os.CreateTemp("", "manifest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write a mock version to the temporary file
	_, err = tempFile.WriteString(`{"version": "1.0.0"}`)
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// Temporarily replace the DefaultManifestFileName with the temporary file's name
	oldDefaultManifestFileName := naming.DefaultManifestFileName
	naming.DefaultManifestFileName = tempFile.Name()
	defer func() { naming.DefaultManifestFileName = oldDefaultManifestFileName }()

	// Call GetCurrentVersion
	version, err := GetCurrentVersion()
	if err != nil {
		t.Fatal(err)
	}

	// Check the returned version
	expectedVersion, _ := semver.NewVersion("1.0.0")
	if !version.Equal(expectedVersion) {
		t.Errorf("Expected version %s, got %s", expectedVersion, version)
	}
}
