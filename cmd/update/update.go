/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/git"
	"github.com/thschue/git-releaser/pkg/helpers"
	"github.com/thschue/git-releaser/pkg/manifest"
	"github.com/thschue/git-releaser/pkg/versioning"
	"os"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var versions config.Versions
		additionalConfig := make(map[string]string)

		if viper.GetString("repository") != "" {
			additionalConfig["repository"] = viper.GetString("repository")
		}

		if viper.GetInt("project_id") != 0 {
			additionalConfig["projectId"] = fmt.Sprintf("%d", viper.GetInt("project_id"))
		}

		g := git.NewGitClient(git.GitConfig{
			Provider:         viper.GetString("provider"),
			AccessToken:      viper.GetString("token"),
			UserId:           viper.GetString("user_id"),
			ProjectUrl:       viper.GetString("project_url"),
			ApiUrl:           viper.GetString("api_url"),
			AdditionalConfig: additionalConfig,
		})

		conf, err := config.ReadConfig(viper.ConfigFileUsed())
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				fmt.Println(err)
			}
		}

		if conf.TargetBranch == "" {
			conf.TargetBranch = "main"
		}

		versions.CurrentVersion, err = manifest.GetCurrentVersion()
		if err != nil {
			fmt.Println("Could not get current version: " + err.Error())
		}

		versions.NextVersion, versions.NewVersion = versioning.GetNextVersion(conf.Versioning)
		versions.VersionPrefix = conf.Versioning.VersionPrefix
		versions.CurrentVersionSlug = versions.VersionPrefix + versions.CurrentVersion.String()
		versions.NextVersionSlug = versions.VersionPrefix + versions.NextVersion.String()
		releaseExists, err := g.CheckRelease(versions)
		if err != nil {
			fmt.Println("Could not check for Release: " + err.Error())
		}

		if !releaseExists {
			fmt.Println("Running release for version " + versions.CurrentVersionSlug)
			err = g.CreateRelease(conf.TargetBranch, versions, "New Release")
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		if !versions.NewVersion {
			fmt.Println("No new version will be created")
			return
		}

		branch, err := g.CheckCreateBranch(conf.TargetBranch, versions.NextVersionSlug)
		if err != nil {
			fmt.Println("Could not check for Branch: " + err.Error())
		}

		content := fmt.Sprintf(`{"version": "%s"}`, versions.NextVersion.String())
		err = g.CommitManifest(branch, content, versions, conf.ExtraFiles)
		if err != nil {
			fmt.Println("Could not update the Repository: " + err.Error())
		}

		err = g.CheckCreatePullRequest(branch, conf.TargetBranch, versions)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	UpdateCmd.Flags().StringP("token", "t", viper.GetString("token"), "Access Token for the Git Provider")
	UpdateCmd.Flags().StringP("api_url", "a", viper.GetString("api_url"), "(optional) API URL for the Git Provider, automatically set for GitHub and GitLab if using the hosted version")
	UpdateCmd.Flags().StringP("project_url", "p", viper.GetString("project_url"), "Project URL for the Git Provider")
	UpdateCmd.Flags().IntP("project_id", "i", viper.GetInt("project_id"), "Project ID when using GitLab")
	UpdateCmd.Flags().StringP("user_id", "u", viper.GetString("user_id"), "User ID")
	UpdateCmd.Flags().StringP("provider", "g", "github", "Git Provider")
	UpdateCmd.Flags().StringP("repository", "r", viper.GetString("repository"), "Repository when using GitHub")
	UpdateCmd.Flags().StringP("target_branch", "b", viper.GetString("target_branch"), "Target Branch (Default: main)")
	helpers.BindViperFlags(UpdateCmd, viper.GetViper())
}
