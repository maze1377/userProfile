package main

import (
	"context"
	"fmt"
	"os"
	"userProfile/pkg/userProfile"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var getAppsCmd = &cobra.Command{
	Use:   "get-user",
	Short: "show user profile information",
	Run:   getApp,
}

func init() {
	rootCmd.AddCommand(getAppsCmd)
}

func getApp(cmd *cobra.Command, clientIDs []string) {
	printVersion()
	ctx := context.Background()
	provider, err := CreateProvider(ctx, cmd)
	if err != nil {
		panicWithError(err, "failed to create provider")
	}

	for _, clientID := range clientIDs {
		fmt.Printf("clientID: %s\n", clientID)
		user, err := provider.GetClientInfo(ctx, &userProfile.ClientInfoRequest{
			ClientID: clientID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't get client info: %v", err)
			continue
		}
		fmt.Println("output:")
		spew.Dump(user)
	}

}
