package versioning

import (
	"github.com/Masterminds/semver"
	"github.com/thenativeweb/get-next-version/conventionalcommits"
	"github.com/thschue/git-releaser/pkg/config"
	"testing"
)

func TestCalculateNextVersion(t *testing.T) {
	versionConfig := config.VersioningConfig{
		BumpPatchMinorPreMajor: true,
		BumpMinorPreMajor:      true,
	}

	currentVersion, _ := semver.NewVersion("1.0.0")

	testCases := []struct {
		name                   string
		conventionalCommitType conventionalcommits.Type
		expectedVersion        string
		expectedHasNextVersion bool
	}{
		{
			name:                   "Chore",
			conventionalCommitType: conventionalcommits.Chore,
			expectedVersion:        "1.0.0",
			expectedHasNextVersion: false,
		},
		{
			name:                   "Fix",
			conventionalCommitType: conventionalcommits.Fix,
			expectedVersion:        "1.0.1",
			expectedHasNextVersion: true,
		},
		{
			name:                   "Feature",
			conventionalCommitType: conventionalcommits.Feature,
			expectedVersion:        "1.0.1",
			expectedHasNextVersion: true,
		},
		{
			name:                   "BreakingChange",
			conventionalCommitType: conventionalcommits.BreakingChange,
			expectedVersion:        "1.1.0",
			expectedHasNextVersion: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			nextVersion, hasNextVersion := calculateNextVersion(versionConfig, currentVersion, []conventionalcommits.Type{testCase.conventionalCommitType})

			if nextVersion.String() != testCase.expectedVersion {
				t.Errorf("Expected version %s, got %s", testCase.expectedVersion, nextVersion.String())
			}

			if hasNextVersion != testCase.expectedHasNextVersion {
				t.Errorf("Expected hasNextVersion %v, got %v", testCase.expectedHasNextVersion, hasNextVersion)
			}
		})
	}
}
