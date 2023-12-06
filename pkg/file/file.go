package file

import (
	"github.com/thschue/git-releaser/pkg/config"
	"os"
	"regexp"
)

func ReplaceVersion(extraFile config.ExtraFileConfig, version string, versionPrefix string) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		return err
	}

	// Define a regular expression to match the version string with the annotation format
	versionRegex := regexp.MustCompile(`(?m)(.*?)(\d+\.\d+\.\d+)(.*?)# x-git-releaser-version`)

	// Replace all occurrences of the version in annotated lines with the new version
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}"+version+"${3}# x-git-releaser-version")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
