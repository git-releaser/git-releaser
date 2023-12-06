/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/pkg/config"
	"github.com/thschue/git-releaser/pkg/git"
	"github.com/thschue/git-releaser/pkg/helpers"
	"github.com/thschue/git-releaser/pkg/manifest"
	"github.com/thschue/git-releaser/pkg/versioning"
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
			fmt.Println("Could not read config file")
		}

		if conf.TargetBranch == "" {
			conf.TargetBranch = "main"
		}

		currentVersion, _ := manifest.GetCurrentVersion()
		nextVersion, isNewVersion := versioning.GetNextVersion(conf.Versioning)

		releaseExists, err := g.CheckRelease(currentVersion.String())
		if err != nil {
			fmt.Println("Could not check for Release: " + err.Error())
		}

		if !releaseExists {
			fmt.Println("Running release for version " + currentVersion.String())
			err = g.CreateRelease(conf.TargetBranch, currentVersion.String(), "New Release")
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		if !isNewVersion {
			fmt.Println("No new version will be created")
			return
		}

		branch, err := g.CheckCreateBranch(conf.TargetBranch, nextVersion.String())
		if err != nil {
			fmt.Println("Could not check for Branch: " + err.Error())
		}

		content := fmt.Sprintf(`{"version": "%s"}`, nextVersion.String())
		err = g.CommitManifest(branch, content, nextVersion.String(), conf.VersionPrefix, conf.ExtraFiles)
		if err != nil {
			fmt.Println("Could not update the Repository: " + err.Error())
		}

		fmt.Println(branch)
		err = g.CheckCreatePullRequest(branch, conf.TargetBranch, currentVersion.String(), nextVersion.String())
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	UpdateCmd.Flags().StringP("token", "t", viper.GetString("token"), "access token")
	UpdateCmd.Flags().StringP("api_url", "a", viper.GetString("api_url"), "api url")
	UpdateCmd.Flags().StringP("project_url", "p", viper.GetString("project_url"), "project url")
	UpdateCmd.Flags().IntP("project_id", "i", viper.GetInt("project_id"), "project id")
	UpdateCmd.Flags().StringP("user_id", "u", viper.GetString("user_id"), "user id")
	UpdateCmd.Flags().StringP("provider", "g", "github", "git provider")
	UpdateCmd.Flags().StringP("repository", "r", viper.GetString("repository"), "github repository")
	UpdateCmd.Flags().StringP("target_branch", "b", viper.GetString("target_branch"), "target branch")
	helpers.BindViperFlags(UpdateCmd, viper.GetViper())
}
