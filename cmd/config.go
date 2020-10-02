package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.NoArgs,
	Short: "Configure global settings",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
