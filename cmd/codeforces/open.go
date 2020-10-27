package codeforces

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/codeforces/open"
	"github.com/cp-tools/cpt/util"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [SPECIFIER]",
	Short: "open required page in default browser",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if given args is a valid specifier.
		problemCnf := util.LoadLocalConf(confSettings)
		if _, err := parseSpecifier(args, problemCnf); err != nil {
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
		// Local (folder) configurations to use.
		problemCnf := util.LoadLocalConf(confSettings)

		modeFlag := cmd.Flags().MustGetString("mode")

		arg, _ := parseSpecifier(args, problemCnf)
		open.Open(arg, modeFlag)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// All flags available to command.
	openCmd.Flags().StringP("mode", "m", "p", "mode to select page to open")

	// All custom completions for command flags.
	openCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"p\tproblem page",
			"d\tdashboard page",
			"s\tsubmission page",
		}
		return modes, cobra.ShellCompDirectiveDefault
	})

}