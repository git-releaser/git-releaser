/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/thschue/git-releaser/cmd"

var (
	version = "dev"
	commit  = "HEAD"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
