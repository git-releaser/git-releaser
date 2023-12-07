package file

import (
	"github.com/thschue/git-releaser/pkg/config"
	"os"
	"regexp"
)

func ReplaceVersionLines(extraFile config.ExtraFileConfig, versions config.Versions) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		return err
	}

	// Define a regular expression to match the version string with the annotation format
	versionRegex := regexp.MustCompile(`(?m)(.*?)(\d+\.\d+\.\d+)(.*?)# x-git-releaser-version`)

	// Replace all occurrences of the version in annotated lines with the new version
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}"+versions.NextVersion.String()+"${3}# x-git-releaser-version")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceVersionBetweenTags(extraFile config.ExtraFileConfig, versions config.Versions) error {
	// Read the contents of the file
	content, err := os.ReadFile(extraFile.Path)
	if err != nil {
		return err
	}

	// Define a regular expression to match the version string and the rest of the text between the start and end tags
	versionRegex := regexp.MustCompile(`(?s)(<!-- x-git-releaser-version-start -->)(.*?)(\d+\.\d+\.\d+)(.*?)(<!-- x-git-releaser-version-end -->)`)

	// Replace the version string with the new version, preserving the rest of the text
	modifiedContent := versionRegex.ReplaceAllString(string(content), "${1}${2}"+versions.NextVersion.String()+"${4}${5}")

	// Write the modified contents back to the file
	err = os.WriteFile(extraFile.Path, []byte(modifiedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
