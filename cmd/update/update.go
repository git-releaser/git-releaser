/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package update

import (
	"errors"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/git"
	"github.com/git-releaser/git-releaser/pkg/helpers"
	"github.com/git-releaser/git-releaser/pkg/versioning"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the repository with the next version and create a release pull request",
	Long:  `Update the repository with the next version and create a release pull request.`,
	Run: func(cmd *cobra.Command, args []string) {
		additionalConfig := make(map[string]string)

		if viper.GetString("repository") != "" {
			additionalConfig["repository"] = viper.GetString("repository")
		}

		if viper.GetInt("project_id") != 0 {
			additionalConfig["projectId"] = fmt.Sprintf("%d", viper.GetInt("project_id"))
		}

		conf, err := config.ReadConfig(viper.ConfigFileUsed())
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				fmt.Println(err)
			}
		}

		g := git.NewGitClient(git.Config{
			Provider:           viper.GetString("provider"),
			AccessToken:        viper.GetString("token"),
			UserId:             viper.GetString("user_id"),
			ProjectUrl:         viper.GetString("project_url"),
			ApiUrl:             viper.GetString("api_url"),
			AdditionalConfig:   additionalConfig,
			PropagationTargets: conf.PropagationTargets,
			DryRun:             viper.GetBool("dry-run"),
			ConfigUpdates:      conf.ConfigUpdates,
		})

		if conf.TargetBranch == "" {
			conf.TargetBranch = "main"
		}

		v := versioning.NewVersion(conf.Versioning)

		err = v.SetNextVersion()
		if err != nil {
			fmt.Println(err)
		}

		versions := v.GetVersions()

		releaseExists, err := g.CheckRelease(versions)
		if err != nil {
			fmt.Println("Could not check for Release: " + err.Error())
		}

		if !releaseExists {
			fmt.Println("Running release for version " + versions.CurrentVersion.Original())
			err = g.CreateRelease(conf.TargetBranch, versions, "")
			if err != nil {
				fmt.Println(err)
			}

			if len(conf.ConfigUpdates) > 0 {
				for _, update := range conf.ConfigUpdates {
					if update.ProjectId != 0 {
						additionalConfig["projectId"] = strconv.Itoa(update.ProjectId)
					}

					r := git.NewGitClient(git.Config{
						Provider:         viper.GetString("provider"),
						AccessToken:      viper.GetString("token"),
						UserId:           viper.GetString("user_id"),
						ProjectUrl:       update.Repository,
						ApiUrl:           viper.GetString("api_url"),
						AdditionalConfig: additionalConfig,
						DryRun:           viper.GetBool("dry-run"),
					})

					changeset, err := r.ReplaceTaggedLines(update.Files, update.SearchTag, versions.CurrentVersion.String())
					if err != nil {
						fmt.Println(err)
					}

					err = r.CommitFile(fmt.Sprintf("release/replace-%s-%s", update.SearchTag, versions.CurrentVersion.String()), changeset)
					if err != nil {
						fmt.Println("Could not update the Repository: " + err.Error())
					}

					err = r.CheckCreateFileMergeRequest(fmt.Sprintf("release/replace-%s-%s", update.SearchTag, versions.CurrentVersion.String()), conf.TargetBranch)
					if err != nil {
						fmt.Println("Could not create the Merge Request: " + err.Error())
					}
				}
			}
			return
		}

		if !versions.HasNextVersion {
			fmt.Println("No new version will be created")
			return
		}

		branch, err := g.CheckCreateBranch(conf.TargetBranch, versions.NextVersion.Original(), conf.BranchPrefix)
		if err != nil {
			fmt.Println("Could not check for Branch: " + err.Error())
		}

		content := fmt.Sprintf(`{"version": "%s"}`, versions.NextVersion.Original())
		err = g.CommitManifest(branch, content, versions, conf.ExtraFiles)
		if err != nil {
			fmt.Println("Could not update the Repository: " + err.Error())
		}

		err = g.CheckCreateReleasePullRequest(branch, conf.TargetBranch, versions)
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
	UpdateCmd.Flags().BoolP("dry-run", "d", viper.GetBool("dry-run"), "Dry-Run")
	helpers.BindViperFlags(UpdateCmd, viper.GetViper())
}
