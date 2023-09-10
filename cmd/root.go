package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	//go:embed embeds/banner.txt
	banner string
	//go:embed embeds/examples.txt
	examples string
)

var rootCmd = &cobra.Command{
	Use:     "API",
	Long:    banner,
	Example: examples,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Init() {}
