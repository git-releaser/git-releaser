/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thschue/git-releaser/cmd/initialize"
	"github.com/thschue/git-releaser/cmd/release"
	"github.com/thschue/git-releaser/cmd/update"
	"os"
)

var (
	cfgFile string
	Version string
	Commit  string
	Date    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-releaser",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string, c string, d string) {
	Version = v
	Commit = c
	Date = d
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-releaser.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(initialize.InitializeCmd)
	rootCmd.AddCommand(update.UpdateCmd)
	rootCmd.AddCommand(release.ReleaseCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Default config file (.git-releaser-config.yaml)")
}

func initConfig() {
	viper.SetConfigName(".git-releaser-config") // name of the config file (without extension)
	viper.AddConfigPath("/etc/myapp/")          // path to look for the config file in
	viper.AddConfigPath("$HOME/.myapp")         // call multiple times to add many search paths
	viper.AddConfigPath(".")                    // look for the config in the working directory
	viper.SetEnvPrefix("GIT_RELEASER")
	viper.AutomaticEnv() // automatically read environment variables

	// Read in config file if found
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found, using default values and/or environment variables")
	}
}
