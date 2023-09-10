package cmd

import (
	"DadJokesAPI/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Starts and runs the API server",
	Long:         banner,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}
