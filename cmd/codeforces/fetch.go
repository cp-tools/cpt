package codeforces

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/codeforces/fetch"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [SPECIFIER]",
	Short: "Fetch and save problem tests",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if given args is a valid specifier.
		problemCnf := util.LoadLocalConf(confSettings)
		if _, err := parseSpecifier(args, problemCnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Load headless browser to use.
		startHeadlessBrowser()

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// Local (folder) configurations to use.
		problemCnf := util.LoadLocalConf(confSettings)

		arg, _ := parseSpecifier(args, problemCnf)
		fetch.Fetch(arg, confSettings)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
