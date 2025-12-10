package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:               "update",
		PersistentPreRunE: createClient,
	}

	updateCfgCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := client.GetConfig()
			if err != nil {
				return fmt.Errorf("getting current config: %w", err)
			}

			if cmd.Flags().Changed("attempt-notify") {
				cfg.AttemptNotify = cfgAttemptNotify
			}

			if cmd.Flags().Changed("max-msg-length") {
				cfg.MaxMessageLength = cfgMaxMsgLen
			}

			if cmd.Flags().Changed("max-concurrent-syncs") {
				cfg.MaxConcurrentSyncs = cfgMaxSyncs
			}

			cfg, err = client.UpdateConfig(cfg)
			if err != nil {
				return err
			}

			fmt.Println("Successfully updated app config. Current config:")
			printCfg(cfg)
			return nil
		},
	}
)

func init() {
	updateCmd.AddCommand(updateCfgCmd)
	updateCfgCmd.Flags().BoolVarP(&cfgAttemptNotify, "attempt-notify", "n", false, "attempt notify on server")
	updateCfgCmd.Flags().IntVarP(&cfgMaxMsgLen, "max-msg-length", "l", 300, "max webex message length")
	updateCfgCmd.Flags().IntVarP(&cfgMaxSyncs, "max-concurrent-syncs", "s", 5, "max concurrent syncs")
}
