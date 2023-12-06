/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package release

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/pkg/git"
	"github.com/thschue/git-releaser/pkg/manifest"

	"github.com/spf13/cobra"
)

var (
	token      string
	apiURL     string
	projectURL string
	projectId  int
	userId     string
	provider   string
	repository string
)

// ReleaseCmd represents the release command
var ReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("TOKEN")
		apiURL := viper.GetString("API_URL")
		projectURL := viper.GetString("PROJECT_URL")
		projectId := viper.GetInt("PROJECT_ID")
		provider := viper.GetString("PROVIDER")
		repository := viper.GetString("REPOSITORY")
		userId := viper.GetString("USER_ID")
		additionalConfig := make(map[string]string)

		if apiURL != "" {
			additionalConfig["apiUrl"] = apiURL
		}

		if repository != "" {
			additionalConfig["repository"] = repository
		}

		if projectId != 0 {
			additionalConfig["projectId"] = fmt.Sprintf("%d", projectId)
		}

		g := git.NewGitClient(git.GitConfig{
			Provider:         provider,
			AccessToken:      token,
			UserId:           userId,
			ProjectUrl:       projectURL,
			AdditionalConfig: additionalConfig,
		})

		currentVersion, err := manifest.GetCurrentVersion()
		if err != nil {
			fmt.Println(err)
		}

		err = g.CreateRelease("main", currentVersion.String(), "New Release")
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	ReleaseCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	ReleaseCmd.Flags().StringVarP(&apiURL, "api-url", "a", "", "api url")
	ReleaseCmd.Flags().StringVarP(&projectURL, "project-url", "p", "", "project url")
	ReleaseCmd.Flags().IntVarP(&projectId, "project-id", "i", 0, "project id")
	ReleaseCmd.Flags().StringVarP(&userId, "user-id", "u", "", "user id")
	ReleaseCmd.Flags().StringVarP(&provider, "provider", "g", "github", "git provider")
	ReleaseCmd.Flags().StringVarP(&repository, "repository", "r", "", "github repository")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
