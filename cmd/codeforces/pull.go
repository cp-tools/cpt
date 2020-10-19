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
		// Check if given args is a valid specifier.
		problemCnf := util.LoadLocalConf(confSettings)
		if _, err := parseSpecifier(args, problemCnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Check if mode flag is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "latest" && modeFlag != "latest-ac" &&
			modeFlag != "all" && modeFlag != "all-ac" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		// Load headless browser to use.
		startHeadlessBrowser()

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Local (folder) configurations to use.
		problemCnf := util.LoadLocalConf(confSettings)

		usernameFlag := cmd.Flags().MustGetString("username")
		modeFlag := cmd.Flags().MustGetString("mode")

		arg, _ := parseSpecifier(args, problemCnf)
		pull.Pull(arg, modeFlag, usernameFlag, confSettings)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// All flags available to command.
	pullCmd.Flags().String("username", "", "user to fetch submissions of")
	pullCmd.Flags().StringP("mode", "m", "latest-ac", "mode to select submissions to save")

	// All custom completions for command flags.
	pullCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"latest",    // Most recent submission.
			"latest-ac", // Most recent AC submission.
			"all",       // All submissions.
			"all-ac",    // All ac submissions.
		}
		return modes, cobra.ShellCompDirectiveDefault
	})
}