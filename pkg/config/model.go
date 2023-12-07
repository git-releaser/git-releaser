package config

import (
	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TargetBranch string            `yaml:"target_branch"`
	Provider     string            `yaml:"provider"`
	ExtraFiles   []ExtraFileConfig `yaml:"extra_files"`
	UserId       string            `yaml:"user_id"`
	AccessToken  string            `yaml:"access_token"`
	ProjectUrl   string            `yaml:"project_url"`
	APIUrl       string            `yaml:"api_url"`
	ProjectID    int               `yaml:"project_id"`
	Repository   string            `yaml:"repository,omitempty"`
	Versioning   VersioningConfig  `yaml:"versioning"`
}

type VersioningConfig struct {
	VersionPrefix          string `yaml:"version_prefix,omitempty"`
	BumpMinorPreMajor      bool   `yaml:"bump_minor_pre_major"`
	BumpPatchMinorPreMajor bool   `yaml:"bump_patch_minor_pre_major"`
}

type ExtraFileConfig struct {
	Path  string `yaml:"path"`
	Label string `yaml:"label,omitempty"`
}

type Versions struct {
	CurrentVersion     *semver.Version `yaml:"current_version"`
	NextVersion        semver.Version  `yaml:"next_version"`
	VersionPrefix      string          `yaml:"version_prefix"`
	NewVersion         bool            `yaml:"new_version"`
	CurrentVersionSlug string          `yaml:"current_version_slug"`
	NextVersionSlug    string          `yaml:"next_version_slug"`
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
