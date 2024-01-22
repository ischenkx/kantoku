package cli

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	root := &cobra.Command{Use: "kanto"}
	root.AddCommand(NewDeploy())
	root.AddCommand(NewEvents())
	root.AddCommand(NewMonitor())

	return root
}
