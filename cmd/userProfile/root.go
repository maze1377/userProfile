package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "userProfile serve",
	Short: "serves profile for the user",
	Long:  "Serves profile from a local database which is accompanied by a cache",
	Run:   nil,
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringP("config-file", "c", "",
		"Path to the config file (eg ./config.yaml) [Optional]")
}
