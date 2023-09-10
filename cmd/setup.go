package cmd

import (
	"DadJokesAPI/setup"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:          "setup",
	Short:        "Configures the server and database",
	Long:         banner,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		setup.Run()
	},
}
