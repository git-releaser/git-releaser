package versioning

import (
	"github.com/Masterminds/semver"
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
}

func NewVersion(cfg config.VersioningConfig) IVersion {
	currentVersion, err := manifest.GetCurrentVersion()
	if err != nil {
		panic(err)
	}

	history, _ := common.GetGitHistory("", currentVersion.Original())

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
