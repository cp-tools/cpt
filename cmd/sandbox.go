// +build !windows

package cmd

import (
	"github.com/cp-tools/cpt/cmd/test"

	"github.com/spf13/cobra"
)

var _sandboxCmd = &cobra.Command{
	Use: "_sandbox",
	Run: func(cmd *cobra.Command, args []string) {
		command := cmd.Flags().MustGetString("command")
		memoryLimit := cmd.Flags().MustGetUint64("memory-limit")

		test.Sandbox(command, memoryLimit)
	},

	Hidden: true,
}

func init() {
	rootCmd.AddCommand(_sandboxCmd)

	// All flags available to command.
	_sandboxCmd.Flags().String("command", "", "")
	_sandboxCmd.Flags().Uint64("memory-limit", 256*1024*1024, "")
}
