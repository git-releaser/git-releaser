/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package update_files

import (
	"errors"
	"fmt"
	"github.com/git-releaser/git-releaser/pkg/config"
	"github.com/git-releaser/git-releaser/pkg/file"
	"github.com/git-releaser/git-releaser/pkg/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var UpdateFilesCmd = &cobra.Command{
	Use:   "update-files",
	Short: "updates tagged lines in files (yaml or json) with new strings",
	Long:  "updates tagged lines in files (yaml or json) with new strings",
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

		if conf.TargetBranch == "" {
			conf.TargetBranch = "main"
		}

		searchString := viper.GetString("search-tag")
		replaceString := viper.GetString("replace-string")
		filePath := viper.GetString("file")

		err = file.ReplaceTaggedLines(filePath, searchString, replaceString)
		if err != nil {
			fmt.Println(err)
		}

		err = file.CommitFile(fmt.Sprintf("release/replace-%s-%s", searchString, replaceString), viper.GetString("user_id"), viper.GetString("token"), filePath, viper.GetBool("dry-run"))
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	UpdateFilesCmd.Flags().StringP("search-tag", "s", viper.GetString("search-tag"), "Tag to search for in the annotation")
	UpdateFilesCmd.Flags().StringP("replace-string", "r", viper.GetString("replace-string"), "String to replace the tag with")
	UpdateFilesCmd.Flags().StringP("file", "f", viper.GetString("file"), "File path to update")
	helpers.BindViperFlags(UpdateFilesCmd, viper.GetViper())
}
