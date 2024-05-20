package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewMonitor() *cobra.Command {
	var monitorCmd = &cobra.Command{
		Use:   "monitor",
		Short: "Monitor the application",
		Run: func(cmd *cobra.Command, args []string) {
			// Placeholder logic for monitor command
			interval, _ := cmd.Flags().GetDuration("interval")
			fmt.Printf("Monitoring with interval: %s\n", interval)
		},
	}

	monitorCmd.Flags().Duration("interval", 100, "Monitoring interval")

	return monitorCmd
}
