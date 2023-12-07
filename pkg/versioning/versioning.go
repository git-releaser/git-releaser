package versioning

import (
	"fmt"
	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/manifest"
	"log"
)

var rootRepositoryFlag string

func GetNextVersion(version config.VersioningConfig) (semver.Version, bool) {
	var nextVersion semver.Version
	var hasNextVersion bool

	repository, err := gogit.PlainOpen(rootRepositoryFlag)
	if err != nil {
		log.Fatal("Could not open repository: " + err.Error())
	}

	result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
	if err != nil {
		log.Fatal("Could not get conventional commits since last release: " + err.Error())
	} else {
		currentVersion, err := manifest.GetCurrentVersion()
		if err != nil {
			log.Fatal("Could not get next version: " + err.Error())
		}
		nextVersion, hasNextVersion = calculateNextVersion(version, currentVersion, result.ConventionalCommitTypes)
	}

	fmt.Println("Next version: ", nextVersion.String())

	return nextVersion, hasNextVersion
}

func calculateNextVersion(
	versionConfig config.VersioningConfig,
	currentVersion *semver.Version,
	conventionalCommitTypes []conventionalcommits.Type,
) (semver.Version, bool) {
	currentlyDetectedChange := conventionalcommits.Chore
	for _, commitType := range conventionalCommitTypes {
		if commitType > currentlyDetectedChange {
			currentlyDetectedChange = commitType
		}
		if currentlyDetectedChange == conventionalcommits.BreakingChange {
			break
		}
	}

	switch currentlyDetectedChange {
	case conventionalcommits.Chore:
		return *currentVersion, false
	case conventionalcommits.Fix:
		return currentVersion.IncPatch(), true
	case conventionalcommits.Feature:
		if versionConfig.BumpPatchMinorPreMajor {
			return currentVersion.IncPatch(), true
		}
		return currentVersion.IncMinor(), true
	case conventionalcommits.BreakingChange:
		if versionConfig.BumpMinorPreMajor {
			return currentVersion.IncMinor(), true
		}
		return currentVersion.IncMajor(), true
	}

	panic("invalid conventional commit type")
}
