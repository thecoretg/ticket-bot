package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thecoretg/ticketbot/internal/cfg"
	"github.com/thecoretg/ticketbot/internal/server"
	"github.com/thecoretg/ticketbot/internal/service"
)

var (
	ctx                                         = context.Background()
	config                                      *cfg.Cfg
	srv                                         *server.Server
	configPath                                  string
	preloadBoards, preloadTickets, initWebhooks bool
	maxPreloads                                 int
	rootCmd                                     = &cobra.Command{
		Use:               "tbot",
		PersistentPreRunE: rootPreRun,
	}

	preloadCmd = &cobra.Command{
		Use:   "preload",
		Short: "Preload boards and/or tickets from Connectwise PSA",
		RunE:  preload,
	}

	hooksCmd = &cobra.Command{
		Use:   "hooks",
		Short: "Initialize webhooks for the TicketBot server in Connectwise PSA",
		RunE:  initHooks,
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the server",
		RunE:  runServer,
	}

	installServiceCmd = &cobra.Command{
		Use:   "install-service",
		Short: "Create a systemd unit for the TicketBot server on a Linux host",
		RunE:  runInstallService,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(preloadCmd)
	rootCmd.AddCommand(hooksCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installServiceCmd)
	preloadCmd.PersistentFlags().BoolVarP(&preloadBoards, "boards", "b", false, "preload boards from connectwise")
	preloadCmd.PersistentFlags().BoolVarP(&preloadTickets, "tickets", "t", false, "preload open tickets from connectwise")
	preloadCmd.PersistentFlags().IntVarP(&maxPreloads, "max-concurrent", "m", 5, "max simultaneous connectwise preloads")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "specify a config file path, otherwise defaults to $HOME/.config/ticketbot")
}

func rootPreRun(cmd *cobra.Command, args []string) error {
	var err error
	config, err = cfg.InitCfg(configPath)
	if err != nil {
		return fmt.Errorf("initializing config: %w", err)
	}

	d, err := server.ConnectToDB(ctx, config.Creds.PostgresDSN)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	srv = server.NewServer(config, d)
	return nil
}

func preload(cmd *cobra.Command, args []string) error {
	if preloadBoards {
		if err := srv.PreloadBoards(ctx, maxPreloads); err != nil {
			return fmt.Errorf("preloading boards: %w", err)
		}
	}

	if preloadTickets {
		if err := srv.PreloadOpenTickets(ctx, maxPreloads); err != nil {
			return fmt.Errorf("preloading tickets: %w", err)
		}
	}

	return nil
}

func initHooks(cmd *cobra.Command, args []string) error {
	if err := srv.InitAllHooks(); err != nil {
		return fmt.Errorf("initializing webhooks: %w", err)
	}

	return nil
}

func runServer(cmd *cobra.Command, args []string) error {
	return srv.Run(ctx)
}

func runInstallService(cmd *cobra.Command, args []string) error {
	if configPath == "" {
		return errors.New("config path is empty, please specify with --config or -c")
	}

	// initialize the config to just to validate it
	_, err := cfg.InitCfg(configPath)
	if err != nil {
		return fmt.Errorf("checking config: %w", err)
	}

	return service.Install(configPath)
}
