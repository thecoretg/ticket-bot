package cli

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiKey  string
	rootCmd = &cobra.Command{
		Use: "tbot",
	}

	webexRoomsCmd = &cobra.Command{
		Use: "rooms",
	}

	cwBoardsCmd = &cobra.Command{
		Use: "boards",
	}

	notifiersCmd = &cobra.Command{
		Use: "notifiers",
	}
)

func Execute() error {
	if apiKey == "" {
		return errors.New("api key is empty")
	}

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func init() {
	apiKey = os.Getenv("TICKETBOT_API_KEY")
}
