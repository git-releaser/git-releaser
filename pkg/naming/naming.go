package naming

import (
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/config"
)

const EnvPrefix = "GIT_RELEASER"
const DefaultConfigFileName = ".git-releaser-config"

var DefaultManifestFileName = ".git-releaser-manifest.json"

func GeneratePrTitle(version string) string {
	title := fmt.Sprintf("Release %s", version)
	return title
}

func CreatePrDescription(version string, changelog string, propagationTargets []config.PropagationTarget, configUpdates []config.ConfigUpdate) string {
	description := fmt.Sprintf("This is a description for the new pull request for version %s.\n\n## Changelog\n\n%s", version, changelog)

	if len(propagationTargets) > 0 {
		description += "\n\n## Propagation Targets\n\n"
		for _, target := range propagationTargets {
			description += fmt.Sprintf("- %s\n", target.Description)
		}
	}

	if len(configUpdates) > 0 {
		description += "\n\n## Configuration Updates\n\n"
		for _, update := range configUpdates {
			description += fmt.Sprintf("- %s - %s\n", update.Repository, update.File)
		}
	}

	return description

}

func CreateReleaseDescription(version string, changelog string) string {
	return fmt.Sprintf("Release %s.\n\n## Changelog\n\n%s", version, changelog)
}

func CreateBranchName(prefix string, version string) string {
	if prefix == "" {
		prefix = "release"
	}
	return fmt.Sprintf("%s-%s", prefix, version)
}
