package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:               "delete",
		PersistentPreRunE: createClient,
	}

	deleteNotifierRuleCmd = &cobra.Command{
		Use:     "notifier-rule",
		Aliases: []string{"rule"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteNotifierRule(id); err != nil {
				return err
			}

			fmt.Printf("Successfully deleted notifier with id of %d\n", id)
			return nil
		},
	}

	deleteForwardCmd = &cobra.Command{
		Use:     "notifier-forward",
		Aliases: []string{"forward", "fwd"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteUserForward(id); err != nil {
				return err
			}

			fmt.Printf("Successfully deleted user forward with id of %d\n", id)
			return nil
		},
	}

	deleteUserCmd = &cobra.Command{
		Use: "user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return errors.New("user id is required")
			}

			if err := client.DeleteUser(id); err != nil {
				return err
			}

			fmt.Printf("User %d successfully deleted\n", id)
			return nil
		},
	}

	deleteAPIKeyCmd = &cobra.Command{
		Use:     "api-key",
		Aliases: []string{"key"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return errors.New("api key id is required")
			}

			if err := client.DeleteAPIKey(id); err != nil {
				return err
			}

			fmt.Printf("API key %d successfully deleted\n", id)
			return nil
		},
	}
)

func init() {
	deleteCmd.AddCommand(deleteNotifierRuleCmd, deleteForwardCmd, deleteUserCmd, deleteAPIKeyCmd)
	deleteForwardCmd.Flags().IntVar(&id, "id", 0, "id of the forward to delete")
	deleteNotifierRuleCmd.Flags().IntVar(&id, "id", 0, "id of the notifier to delete")
	deleteAPIKeyCmd.Flags().IntVar(&id, "id", 0, "id of the key to delete")
	deleteUserCmd.Flags().IntVar(&id, "id", 0, "id of the user to delete")
}
