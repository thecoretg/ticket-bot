package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thecoretg/ticketbot/internal/cfg"
	"github.com/thecoretg/ticketbot/internal/server"
)

var (
	ctx                                              = context.Background()
	run, preloadBoards, preloadTickets, initWebhooks bool
	maxPreloads                                      int
	rootCmd                                          = &cobra.Command{
		Use: "ticketbot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return parseRootFlags(ctx)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&preloadBoards, "preload-boards", "b", false, "preload boards from connectwise")
	rootCmd.PersistentFlags().BoolVarP(&preloadTickets, "preload-tickets", "t", false, "preload open tickets from connectwise")
	rootCmd.PersistentFlags().IntVarP(&maxPreloads, "max-preloads", "m", 5, "max simultaneous connectwise preloads")
	rootCmd.PersistentFlags().BoolVarP(&initWebhooks, "init-webhooks", "w", false, "initialize webhooks")
	rootCmd.PersistentFlags().BoolVarP(&run, "run", "r", false, "run the server")
}

func parseRootFlags(ctx context.Context) error {
	c, err := cfg.InitCfg()
	if err != nil {
		return fmt.Errorf("initializing config: %w", err)
	}

	d, err := server.ConnectToDB(ctx, c.Creds.PostgresDSN)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	s := server.NewServer(c, d)

	if err := s.BootstrapAdmin(ctx); err != nil {
		return fmt.Errorf("bootstrapping admin: %w", err)
	}

	if preloadBoards {
		if err := s.PreloadBoards(ctx, maxPreloads); err != nil {
			return fmt.Errorf("preloading boards: %w", err)
		}
	}

	if preloadTickets {
		if err := s.PreloadOpenTickets(ctx, maxPreloads); err != nil {
			return fmt.Errorf("preloading tickets: %w", err)
		}
	}

	if initWebhooks {
		if err := s.InitAllHooks(); err != nil {
			return fmt.Errorf("initializing webhooks: %w", err)
		}
	}

	if run {
		return s.Run()
	}

	return nil
}
