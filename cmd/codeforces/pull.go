package codeforces

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/codeforces/pull"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [SPECIFIER]",
	Short: "Pulls submissions to local storage",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cnf = util.LoadLocalConf(cnf)

		// Check if given args is a valid specifier.
		if _, err := parseSpecifier(args, cnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Check if mode flag is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "l" && modeFlag != "la" &&
			modeFlag != "a" && modeFlag != "aa" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		// Load headless browser to use.
		startHeadlessBrowser()

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		usernameFlag := cmd.Flags().MustGetString("username")
		modeFlag := cmd.Flags().MustGetString("mode")

		arg, _ := parseSpecifier(args, cnf)
		pull.Pull(arg, modeFlag, usernameFlag, cnf)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// All flags available to command.
	pullCmd.Flags().String("username", "", "user to fetch submissions of")
	pullCmd.Flags().StringP("mode", "m", "la", "mode to select submissions to save")

	// All custom completions for command flags.
	pullCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"l\tlatest",     // Most recent submission.
			"la\tlatest ac", // Most recent AC submission.
			"a\tall",        // All submissions.
			"aa\tall ac",    // All ac submissions.
		}
		return modes, cobra.ShellCompDirectiveDefault
	})
}
