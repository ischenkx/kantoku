package cli

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"os"
)

func New() *cobra.Command {
	root := &cobra.Command{
		Use: "kanto",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Read environment variables from the file provided by the --env flag
			envFilePath, _ := cmd.Flags().GetString("env")
			if envFilePath != "" {
				if err := godotenv.Load(envFilePath); err != nil {
					cmd.PrintErrf("error reading environment variables: %s\n", err)
					os.Exit(1)
				}
			}
		},
	}

	root.PersistentFlags().String("env", "", "Path to a file with environment variables")

	root.AddCommand(NewDeploy())
	root.AddCommand(NewEvents())
	root.AddCommand(NewMonitor())

	return root
}
