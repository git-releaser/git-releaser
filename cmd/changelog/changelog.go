/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package changelog

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/pkg/changelog"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/git"
	"github.com/thschue/git-releaser/pkg/helpers"
	"os"
)

var ChangeLogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		sinceVersion := viper.GetString("since_version")

		if sinceVersion == "" {
			sinceVersion, err = g.GetHighestRelease()
			fmt.Println("sinceVersion: " + conf.Versioning.VersionPrefix + sinceVersion)
			if err != nil {
				fmt.Println(err)
			}
		}

		commits, err := g.GetCommitsSinceRelease(conf.Versioning.VersionPrefix + sinceVersion)
		if err != nil {
			fmt.Println(err)
		}
		conventionalCommits := changelog.ParseConventionalCommits(commits)
		log := changelog.GenerateChangelog(conventionalCommits, viper.GetString("project_url"))
		fmt.Println("Last Version: " + viper.GetString("since_version"))
		fmt.Println(log)
	},
}

func init() {
	ChangeLogCmd.Flags().StringP("token", "t", viper.GetString("token"), "access token")
	ChangeLogCmd.Flags().StringP("api_url", "a", viper.GetString("api_url"), "api url")
	ChangeLogCmd.Flags().StringP("project_url", "p", viper.GetString("project_url"), "project url")
	ChangeLogCmd.Flags().IntP("project_id", "i", viper.GetInt("project_id"), "project id")
	ChangeLogCmd.Flags().StringP("user_id", "u", viper.GetString("user_id"), "user id")
	ChangeLogCmd.Flags().StringP("provider", "g", "github", "git provider")
	ChangeLogCmd.Flags().StringP("repository", "r", viper.GetString("repository"), "github repository")
	ChangeLogCmd.Flags().StringP("target_branch", "b", viper.GetString("target_branch"), "target branch")
	ChangeLogCmd.Flags().StringP("since_version", "l", viper.GetString("since_version"), "version")
	helpers.BindViperFlags(ChangeLogCmd, viper.GetViper())
}
