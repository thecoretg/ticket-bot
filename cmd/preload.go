package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	preloadCmd = &cobra.Command{
		Use: "preload",
	}

	preloadBoardsCmd = &cobra.Command{
		Use:  "boards",
		RunE: preloadBoards,
	}

	preloadTicketsCmd = &cobra.Command{
		Use:  "tickets",
		RunE: preloadTickets,
	}

	preloadAllCmd = &cobra.Command{
		Use:  "all",
		RunE: preloadAll,
	}
)

func addPreloadCmd() {
	rootCmd.AddCommand(preloadCmd)
	preloadCmd.AddCommand(preloadBoardsCmd, preloadTicketsCmd, preloadAllCmd)

	preloadCmd.PersistentFlags().IntVarP(&maxPreloads, "max-concurrent", "m", 5, "max simultaneous connectwise preloads")
}

func preloadBoards(cmd *cobra.Command, args []string) error {
	if err := srv.PreloadBoards(ctx, maxPreloads); err != nil {
		return fmt.Errorf("preloading boards: %w", err)
	}

	return nil
}

func preloadTickets(cmd *cobra.Command, args []string) error {
	if err := srv.PreloadOpenTickets(ctx, maxPreloads); err != nil {
		return fmt.Errorf("preloading tickets: %w", err)
	}

	return nil
}

func preloadAll(cmd *cobra.Command, args []string) error {
	if err := srv.PreloadBoards(ctx, maxPreloads); err != nil {
		return fmt.Errorf("preloading boards: %w", err)
	}

	if err := srv.PreloadOpenTickets(ctx, maxPreloads); err != nil {
		return fmt.Errorf("preloading tickets: %w", err)
	}

	return nil
}
