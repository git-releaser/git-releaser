package naming

import (
	"fmt"
)

const EnvPrefix = "GIT_RELEASER"
const DefaultConfigFileName = ".git-releaser-config"

var DefaultManifestFileName = ".git-releaser-manifest.json"

func GeneratePrTitle(version string) string {
	title := fmt.Sprintf("Release %s", version)
	return title
}

func CreatePrDescription(version string, changelog string) string {
	return fmt.Sprintf("This is a description for the new pull request for version %s.\n\n## Changelog\n\n%s", version, changelog)
}

func CreateReleaseDescription(version string, changelog string) string {
	return fmt.Sprintf("Release %s.\n\n## Changelog\n\n%s", version, changelog)
}
