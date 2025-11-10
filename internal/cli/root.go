package cli

import (
	"github.com/spf13/cobra"
	"github.com/thecoretg/ticketbot/internal/sdk"
)

var (
	boardID int
	roomID  int

	client  *sdk.Client
	rootCmd = &cobra.Command{
		Use:               "tbot",
		PersistentPreRunE: createClient,
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(pingCmd, authCheckCmd, webexRoomsCmd, cwBoardsCmd, notifiersCmd)
}
