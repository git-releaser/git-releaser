/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package initialize

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/pkg/helpers"
	"github.com/thschue/git-releaser/pkg/naming"
	"os"

	"github.com/spf13/cobra"
)

const ReleaseManifestFilename = ".git-releaser-manifest.json"

// initializeCmd represents the initialize command
var InitializeCmd = &cobra.Command{
	Use:   "initialize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(ReleaseManifestFilename); errors.Is(err, os.ErrNotExist) {
			version := "0.0.0"
			// write version to .git-releaser-manifest.json
			manifest, err := os.Create(ReleaseManifestFilename)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer manifest.Close()
			_, err = manifest.WriteString(fmt.Sprintf(`{"version": "%s"}`, version))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println(ReleaseManifestFilename + " already exists")
		}

		if viper.ConfigFileUsed() == "" {
			// write config to .git-releaser.yaml
			err := viper.WriteConfigAs(naming.DefaultConfigFileName + ".yaml")
			if err != nil {
				return
			}
		}
	},
}

func init() {
	InitializeCmd.Flags().StringP("provider", "g", "github", "Git Provider")
	helpers.BindViperFlags(InitializeCmd, viper.GetViper())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initializeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initializeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
