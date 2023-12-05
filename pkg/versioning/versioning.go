package versioning

import (
	"fmt"
	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/versioning"
	"github.com/thschue/git-releaser/pkg/manifest"
	"log"
)

var rootRepositoryFlag string

func GetNextVersion() (semver.Version, bool) {
	var nextVersion semver.Version
	var hasNextVersion bool

	repository, err := gogit.PlainOpen(rootRepositoryFlag)
	if err != nil {
		log.Fatal(err)
	}

	result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
	if err != nil {
		log.Fatal(err)
	} else {
		currentVersion, err := manifest.GetCurrentVersion()
		if err != nil {
			log.Fatal(err)
		}
		nextVersion, hasNextVersion = versioning.CalculateNextVersion(currentVersion, result.ConventionalCommitTypes)
	}

	fmt.Println("Next version: ", nextVersion)
	fmt.Println("Has next version: ", hasNextVersion)

	return nextVersion, hasNextVersion
}
