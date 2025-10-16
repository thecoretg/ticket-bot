package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thecoretg/ticketbot/internal/service"
)

var (
	serviceCmd = &cobra.Command{
		Use: "service",
	}

	serviceInstallCmd = &cobra.Command{
		Use:  "install",
		RunE: runInstallService,
	}

	serviceLogCmd = &cobra.Command{
		Use:  "logs",
		RunE: showServiceLogs,
	}
)

func runInstallService(cmd *cobra.Command, args []string) error {
	return service.Install(configPath)
}

func showServiceLogs(cmd *cobra.Command, args []string) error {
	return service.ShowLogs()
}

func addServiceCmd() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceInstallCmd, serviceLogCmd)
	serviceCmd.AddCommand(serviceLogCmd)
}
