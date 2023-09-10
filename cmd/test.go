package cmd

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:          "test",
	Short:        "Runs unit tests for the API",
	Long:         banner,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTests()
	},
}

func runTests() error {
	command := exec.Command("go", "test", "./tests")
	stdout, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	scanner := bufio.NewScanner(stdout)

	command.Start()
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	command.Wait()

	return nil
}
