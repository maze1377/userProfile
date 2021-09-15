package main

import (
	"fmt"

	"userProfile/pkg/userProfile"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info",
	Long:  `All software have versions. This is aggregator`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("%-18s %-18s Commit:%s                  (%s)\n host: %s StartTime: %s\n", userProfile.Title, userProfile.Version,
		userProfile.Commit, userProfile.BuildTime, userProfile.Hostname, userProfile.StartTime)
}
