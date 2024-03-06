package conventional

import (
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thschue/git-releaser/pkg/config"
)

type ChangeType int

const (
	Major ChangeType = iota
	Minor
	Patch
)

type ConventionalCommitTypesResult struct {
	ConventionalCommitTypes []conventionalcommits.Type
}

type Version struct {
	config.Versions
}

func (v *Version) GetCurrentVersion() semver.Version {
	return v.CurrentVersion
}

func (v *Version) GetNextVersion() (semver.Version, bool) {
	return v.NextVersion, v.HasNextVersion
}

func (v *Version) SetNextVersion() error {
	commitTypes, err := v.getConventionalCommitTypes()
	if err != nil {
		return err
	}

	v.NextVersion, v.HasNextVersion = v.calculateNextVersion(commitTypes.ConventionalCommitTypes)

	return err
}

func (v *Version) GetVersions() config.Versions {
	return v.Versions
}

func (v *Version) GetHistory() []object.Commit { return v.Commits }

func (v *Version) getConventionalCommitTypes() (ConventionalCommitTypesResult, error) {
	conventionalCommitTypes := []conventionalcommits.Type{}

	for _, commit := range v.Commits {
		currentCommitType, err := conventionalcommits.CommitMessageToType(commit.Message)
		if err != nil {
			currentCommitType = conventionalcommits.Chore
		}
		conventionalCommitTypes = append(
			conventionalCommitTypes,
			currentCommitType,
		)
	}

	return ConventionalCommitTypesResult{
		ConventionalCommitTypes: conventionalCommitTypes,
	}, nil
}

func (v *Version) calculateNextVersion(conventionalCommitTypes []conventionalcommits.Type) (semver.Version, bool) {
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
	case conventionalcommits.Fix:
		if v.Config.BumpPatchMinorPreMajor {
			return v.CurrentVersion.IncPatch(), true
		}
		return v.CurrentVersion.IncPatch(), true
	case conventionalcommits.Feature:
		if v.Config.BumpMinorPreMajor {
			return v.CurrentVersion.IncPatch(), true
		}
		return v.CurrentVersion.IncMinor(), true
	case conventionalcommits.BreakingChange:
		return v.CurrentVersion.IncMajor(), true
	default:
		return v.CurrentVersion, false
	}
}
