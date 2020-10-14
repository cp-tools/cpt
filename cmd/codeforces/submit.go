package codeforces

import (
	"fmt"
	"strings"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [SPECIFIER]",
	Short: "Submit problem solution to judge",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load headless browser to use.
		startHeadlessBrowser()

		// Check if given args is a valid specifier.
		if _, err := codeforces.Parse(strings.Join(args, "")); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// All flags available to command.
	submitCmd.Flags().StringP("file", "f", "", "code file to submit")
}
