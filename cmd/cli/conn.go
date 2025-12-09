package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	pingCmd = &cobra.Command{
		Use:               "ping",
		PersistentPreRunE: createClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.Ping(); err != nil {
				return err
			}

			fmt.Println("Server is up")
			return nil
		},
	}

	authCheckCmd = &cobra.Command{
		Use:               "authcheck",
		PersistentPreRunE: createClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.AuthTest(); err != nil {
				return err
			}

			fmt.Println("Successfully authenticated")
			return nil
		},
	}
)
