package codeforces

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/codeforces/open"
	"github.com/cp-tools/cpt/pkg/conf"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [SPECIFIER]",
	Short: "Open required page in default browser",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cnf = conf.New("local").SetParent(cnf).LoadFile("meta.yaml")

		// Check if given args is a valid specifier.
		if _, err := parseSpecifier(args, cnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Check if mode flag is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "p" && modeFlag != "d" && modeFlag != "s" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		mode := cmd.Flags().MustGetString("mode")

		arg, _ := parseSpecifier(args, cnf)
		open.Open(arg, mode)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// All flags available to command.
	openCmd.Flags().StringP("mode", "m", "p", "mode to select page to open")

	// All custom completions for command flags.
	openCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"p\tproblem",
			"d\tdashboard",
			"s\tsubmission",
		}
		return modes, cobra.ShellCompDirectiveDefault
	})
}
