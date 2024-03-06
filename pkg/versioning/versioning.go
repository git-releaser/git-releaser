package versioning

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/git/common"
	"github.com/thschue/git-releaser/pkg/manifest"
	"github.com/thschue/git-releaser/pkg/versioning/conventional"
	"github.com/thschue/git-releaser/pkg/versioning/simple"
)

type IVersion interface {
	GetCurrentVersion() semver.Version
	SetNextVersion() error
	GetNextVersion() (semver.Version, bool)
	GetVersions() config.Versions
	GetHistory() []object.Commit
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
