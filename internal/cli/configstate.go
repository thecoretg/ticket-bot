package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thecoretg/ticketbot/internal/server"
)

var (
	cfgDebug           bool
	attemptNotify      bool
	maxMsgLength       int
	maxConcurrentSyncs int

	stateCmd = &cobra.Command{
		Use: "state",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := client.GetAppState()
			if err != nil {
				return fmt.Errorf("getting app state: %w", err)
			}

			fmt.Printf("Syncing Tickets: %v\nSyncing Rooms: %v\n",
				state.SyncingTickets, state.SyncingWebexRooms)

			return nil
		},
	}

	cfgCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
	}

	getCfgCmd = &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := client.GetConfig()
			if err != nil {
				return fmt.Errorf("getting current config: %w", err)
			}

			printCfg(cfg)
			return nil
		},
	}

	updateCfgCmd = &cobra.Command{
		Use: "update",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := &server.AppConfigPayload{}
			if cmd.Flags().Changed("debug") {
				fmt.Printf("Changing debug to: %v\n", cfgDebug)
				p.Debug = &cfgDebug
			}

			if cmd.Flags().Changed("attempt-notify") {
				fmt.Printf("Changing attempt notify to: %v\n", attemptNotify)
				p.AttemptNotify = &attemptNotify
			}

			if cmd.Flags().Changed("max-msg-length") {
				fmt.Printf("Changing max message length to: %d\n", maxMsgLength)
				p.MaxMessageLength = &maxMsgLength
			}

			if cmd.Flags().Changed("max-concurrent-syncs") {
				fmt.Printf("Changing max concurrent syncs to: %d\n", maxConcurrentSyncs)
				p.MaxConcurrentSyncs = &maxConcurrentSyncs
			}

			cfg, err := client.UpdateConfig(p)
			if err != nil {
				return fmt.Errorf("updating config: %w", err)
			}

			fmt.Println("Successfully updated app config. Current config:")
			printCfg(cfg)
			return nil
		},
	}
)

func printCfg(cfg *server.AppConfig) {
	fmt.Printf("Debug: %v\nAttempt Notify: %v\nMax Msg Length: %d\nMax Concurrent Syncs: %d\n",
		cfg.Debug, cfg.AttemptNotify, cfg.MaxMessageLength, cfg.MaxConcurrentSyncs)
}

func init() {
	cfgCmd.AddCommand(getCfgCmd, updateCfgCmd)
	updateCfgCmd.Flags().BoolVarP(&cfgDebug, "debug", "d", false, "enable debug mode on server")
	updateCfgCmd.Flags().BoolVarP(&attemptNotify, "attempt-notify", "n", false, "attempt notify on server")
	updateCfgCmd.Flags().IntVarP(&maxMsgLength, "max-msg-length", "l", 300, "max webex message length")
	updateCfgCmd.Flags().IntVarP(&maxConcurrentSyncs, "max-concurrent-syncs", "s", 5, "max concurrent syncs")
}
