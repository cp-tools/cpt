package codeforces

import (
	"fmt"
	"strings"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/cmd/codeforces/fetch"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [SPECIFIER]",
	Short: "Fetch and save problem tests",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Load headless browser to use.
		startHeadlessBrowser()

		// Check if given args is a valid specifier.
		if _, err := codeforces.Parse(strings.Join(args, "")); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// Parse args to specifier.
		arg, _ := codeforces.Parse(strings.Join(args, ""))
		fetch.Fetch(arg, confSettings)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
