package codeforces

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/codeforces/fetch"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [SPECIFIER]",
	Short: "Fetch and save problem tests",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cnf = conf.New("local").SetParent(cnf).LoadFile("meta.yaml")

		if _, err := parseSpecifier(args, cnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Load headless browser to use.
		startHeadlessBrowser()

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		arg, _ := parseSpecifier(args, cnf)
		fetch.Fetch(arg, cnf)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
