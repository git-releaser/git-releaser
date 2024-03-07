package config

import (
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TargetBranch       string              `yaml:"target_branch"`
	Provider           string              `yaml:"provider"`
	ExtraFiles         []ExtraFileConfig   `yaml:"extra_files"`
	UserId             string              `yaml:"user_id"`
	AccessToken        string              `yaml:"access_token"`
	ProjectUrl         string              `yaml:"project_url"`
	APIUrl             string              `yaml:"api_url"`
	ProjectID          int                 `yaml:"project_id"`
	Repository         string              `yaml:"repository,omitempty"`
	Versioning         VersioningConfig    `yaml:"versioning"`
	PropagationTargets []PropagationTarget `yaml:"propagation_targets,omitempty"`
}

type VersioningConfig struct {
	VersionPrefix          string            `yaml:"version_prefix,omitempty"`
	BumpMinorPreMajor      bool              `yaml:"bump_minor_pre_major"`
	BumpPatchMinorPreMajor bool              `yaml:"bump_patch_minor_pre_major"`
	Strategy               string            `yaml:"strategy"`
	SimpleCommitTypes      SimpleCommitTypes `yaml:"simple_commit_types,omitempty"`
}

type SimpleCommitTypes struct {
	Patch        []string `yaml:"patch"`
	Minor        []string `yaml:"minor"`
	Major        []string `yaml:"major"`
	DefaultPatch bool     `yaml:"default_patch"`
}

type PropagationTarget struct {
	TargetBranch string `yaml:"target_branch"`
	Target       string `yaml:"target"`
}

type ExtraFileConfig struct {
	Path  string `yaml:"path"`
	Label string `yaml:"label,omitempty"`
}

type Versions struct {
	CurrentVersion semver.Version
	Commits        []object.Commit
	Config         VersioningConfig
	NextVersion    semver.Version
	HasNextVersion bool
}

func ReadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
