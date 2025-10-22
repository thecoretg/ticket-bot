package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	keyDir       string
	bootstrapCmd = &cobra.Command{
		Use: "bootstrap",
	}

	bootstrapCreateCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := srv.BootstrapAdmin(ctx, keyDir)
			if err != nil {
				return err
			}

			fmt.Printf("Bootstrap key created and is stored at %s\n", path)
			return nil
		},
	}
)

func addBootstrapCmd() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.AddCommand(bootstrapCreateCmd)
	bootstrapCreateCmd.Flags().StringVarP(&keyDir, "key-directory", "k", "", "directory to put the bootstrap key in")
}
