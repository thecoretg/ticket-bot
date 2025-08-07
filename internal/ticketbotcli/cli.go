package ticketbotcli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"tctg-automation/internal/ticketbot"
)

var (
	ctx            context.Context
	server         *ticketbot.Server
	preloadBoards  bool
	preloadTickets bool

	rootCmd = &cobra.Command{
		Use: "ticketbot",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	adminCmd = &cobra.Command{
		Use: "admin",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx = context.Background()
			cfg, err := ticketbot.InitCfg(ctx)
			if err != nil {
				return fmt.Errorf("an error occured initializing the config: %w", err)
			}

			server, err = ticketbot.NewServer(cfg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	preloadCmd = &cobra.Command{
		Use: "preload",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.PreloadAll(ctx, preloadBoards, preloadTickets)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	preloadCmd.PersistentFlags().BoolVarP(&preloadBoards, "boards", "b", false, "set to preload boards")
	preloadCmd.PersistentFlags().BoolVarP(&preloadTickets, "tickets", "p", false, "set to preload tickets")

	rootCmd.AddCommand(adminCmd)
	adminCmd.AddCommand(preloadCmd)
}
