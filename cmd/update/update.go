/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"fmt"
	"github.com/Masterminds/semver"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/thenativeweb/get-next-version/git"
	"github.com/thenativeweb/get-next-version/versioning"
	"github.com/thschue/git-releaser/pkg/manifest"
	"log"
)

var rootRepositoryFlag string
var rootTargetFlag string
var rootPrefixFlag string

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
		getNextVersion()
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getNextVersion() {
	var nextVersion semver.Version
	var hasNextVersion bool

	repository, err := gogit.PlainOpen(rootRepositoryFlag)
	if err != nil {
		log.Fatal(err)
	}

	result, err := git.GetConventionalCommitTypesSinceLastRelease(repository)
	if err != nil {
		log.Fatal(err)
	} else {
		currentVersion, err := manifest.GetCurrentVersion()
		if err != nil {
			log.Fatal(err)
		}
		nextVersion, hasNextVersion = versioning.CalculateNextVersion(currentVersion, result.ConventionalCommitTypes)
	}

	fmt.Println(nextVersion)
	fmt.Println(hasNextVersion)
}
