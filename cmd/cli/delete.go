package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use: "delete",
	}

	deleteNotifierRuleCmd = &cobra.Command{
		Use:     "notifier-rule",
		Aliases: []string{"rule"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteNotifierRule(notifierID); err != nil {
				return err
			}

			fmt.Printf("Successfully deleted notifier with id of %d\n", notifierID)
			return nil
		},
	}

	deleteForwardCmd = &cobra.Command{
		Use:     "notifier-forward",
		Aliases: []string{"forward", "fwd"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteUserForward(forwardID); err != nil {
				return err
			}

			fmt.Printf("Successfully deleted user forward with id of %d\n", forwardID)
			return nil
		},
	}
)

func init() {
	deleteCmd.AddCommand(deleteNotifierRuleCmd, deleteForwardCmd)
	deleteNotifierRuleCmd.Flags().IntVar(&id, "id", 0, "id of the rule to delete")
	deleteForwardCmd.Flags().IntVar(&id, "id", 0, "id of the forward to delete")
}
