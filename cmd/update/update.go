/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/pkg/git"
	"github.com/thschue/git-releaser/pkg/versioning"
)

var (
	token      string
	apiURL     string
	projectURL string
	projectId  int
	userId     string
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
		token := viper.GetString("TOKEN")
		apiURL := viper.GetString("API_URL")
		projectURL := viper.GetString("PROJECT_URL")
		projectId := viper.GetInt("PROJECT_ID")

		additionalConfig := make(map[string]string)

		if apiURL != "" {
			additionalConfig["apiUrl"] = apiURL
		}
		if projectId != 0 {
			additionalConfig["projectId"] = fmt.Sprintf("%d", projectId)
		}

		g := git.NewGitClient(git.GitConfig{
			AccessToken:      token,
			UserId:           userId,
			ProjectUrl:       projectURL,
			AdditionalConfig: additionalConfig,
		})

		branch, err := g.CheckCreateBranch()
		if err != nil {
			fmt.Println(err)
		}

		nextVersion, _ := versioning.GetNextVersion()
		content := fmt.Sprintf(`{"version": "%s"}`, nextVersion.String())
		err = g.CommitManifest(branch, content)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(branch)
		err = g.CheckCreatePullRequest(branch, "main")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	UpdateCmd.Flags().StringVarP(&apiURL, "api-url", "a", "", "api url")
	UpdateCmd.Flags().StringVarP(&projectURL, "project-url", "p", "", "project url")
	UpdateCmd.Flags().IntVarP(&projectId, "project-id", "i", 0, "project id")
	UpdateCmd.Flags().StringVarP(&userId, "user-id", "u", "", "user id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
