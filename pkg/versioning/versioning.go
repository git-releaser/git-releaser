package versioning

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/git/common"
	"github.com/git-releaser/git-releaser/pkg/manifest"
	"github.com/git-releaser/git-releaser/pkg/versioning/conventional"
	"github.com/git-releaser/git-releaser/pkg/versioning/simple"
)

type IVersion interface {
	GetCurrentVersion() semver.Version
	SetNextVersion() error
	GetNextVersion() (semver.Version, bool)
	GetVersions() config.Versions
}

func NewVersion(cfg config.VersioningConfig) IVersion {
	currentVersion, err := manifest.GetCurrentVersion()
	if err != nil {
		panic(err)
	}

	history, err := common.GetGitHistory("", currentVersion.Original())
	if err != nil {
		fmt.Println("Could not get git history:" + err.Error())
	}

	switch cfg.Strategy {
	case "conventional":
		return &conventional.Version{
			Versions: config.Versions{
				CurrentVersion: *currentVersion,
				Commits:        history,
				HasNextVersion: false,
				Config:         cfg,
			},
		}
	default:
		return &simple.Version{
			Versions: config.Versions{
				CurrentVersion: *currentVersion,
				Commits:        history,
				HasNextVersion: false,
				Config:         cfg,
			},
		}
	}
}
