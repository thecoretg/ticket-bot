package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	enableNotify bool
	notifierID   int

	notifiersCmd = &cobra.Command{
		Use: "notifiers",
	}

	listNotifiersCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, err := client.ListNotifiers()
			if err != nil {
				return fmt.Errorf("retrieving notifiers: %w", err)
			}

			if ns == nil || len(ns) == 0 {
				fmt.Println("No notifiers found")
				return nil
			}

			notifiersTable(ns)
			return nil
		},
	}

	getNotifierCmd = &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := client.GetNotifier(notifierID)
			if err != nil {
				return fmt.Errorf("retrieving notifier: %w", err)
			}

			fmt.Printf("ID: %d\nRoom: %s\nBoard: %s\nNotify: %v\n",
				n.ID, n.WebexRoom.Name, n.CWBoard.Name, n.NotifyEnabled)

			return nil
		},
	}

	createNotifierCmd = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := client.CreateNotifier(boardID, roomID, enableNotify)
			if err != nil {
				return fmt.Errorf("creating notifier: %w", err)
			}

			fmt.Printf("ID: %d\nRoom: %s\nBoard: %s\nNotify: %v\n",
				n.ID, n.WebexRoom.Name, n.CWBoard.Name, n.NotifyEnabled)

			return nil
		},
	}

	deleteNotifierCmd = &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteNotifier(notifierID); err != nil {
				return fmt.Errorf("deleting notifier: %w", err)
			}

			fmt.Printf("Successfully deleted notifier with id of %d\n", notifierID)
			return nil
		},
	}
)

func init() {
	notifiersCmd.AddCommand(createNotifierCmd, listNotifiersCmd, getNotifierCmd, deleteNotifierCmd)
	notifiersCmd.PersistentFlags().IntVar(&notifierID, "id", 0, "id of notifier")
	createNotifierCmd.Flags().IntVarP(&boardID, "board-id", "b", 0, "board id to use")
	createNotifierCmd.Flags().IntVarP(&roomID, "room-id", "r", 0, "room id to use")
	createNotifierCmd.Flags().BoolVarP(&enableNotify, "enable-notify", "n", false, "enable notify for rule")
}
