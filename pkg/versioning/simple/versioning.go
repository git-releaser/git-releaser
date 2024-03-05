package simple

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/thschue/git-releaser/pkg/config"
	"strings"
)

type ChangeType int

const (
	Major ChangeType = iota
	Minor
	Patch
	None
)

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
	commitTypes := v.getChangeTypes()

	if len(commitTypes) == 0 {
		return fmt.Errorf("no change types found")
	}

	v.NextVersion, v.HasNextVersion = v.calculateNextVersion(commitTypes)

	return nil
}

func (v *Version) calculateNextVersion(changeTypes []ChangeType) (semver.Version, bool) {
	currentlyDetectedChange := None
	for _, commitType := range changeTypes {
		if commitType < currentlyDetectedChange {
			currentlyDetectedChange = commitType
		}
		if currentlyDetectedChange == Major {
			break
		}
	}

	switch currentlyDetectedChange {
	case Patch:
		fmt.Println(currentlyDetectedChange)
		return v.CurrentVersion.IncPatch(), true
	case Minor:
		fmt.Println(currentlyDetectedChange)
		return v.CurrentVersion.IncMinor(), true
	case Major:
		fmt.Println(currentlyDetectedChange)
		return v.CurrentVersion.IncMajor(), true
	default:
		if v.Versions.Config.SimpleCommitTypes.DefaultPatch {
			return v.CurrentVersion.IncPatch(), true
		}
		return v.CurrentVersion, false
	}
}

func (v *Version) getChangeTypes() []ChangeType {
	changeTypes := []ChangeType{}

	for _, commit := range v.Commits {
		fmt.Println(commit.Message)
		for _, prefix := range v.Config.SimpleCommitTypes.Major {
			if strings.HasPrefix(commit.Message, prefix) {
				changeTypes = append(changeTypes, Major)
			}
		}
		for _, prefix := range v.Config.SimpleCommitTypes.Minor {
			if strings.HasPrefix(commit.Message, prefix) {
				changeTypes = append(changeTypes, Minor)
			}
		}
		for _, prefix := range v.Config.SimpleCommitTypes.Patch {
			if strings.HasPrefix(commit.Message, prefix) {
				changeTypes = append(changeTypes, Patch)
			}
		}
		changeTypes = append(changeTypes, None)
	}
	return changeTypes
}

func (v *Version) GetVersions() config.Versions {
	return v.Versions
}
