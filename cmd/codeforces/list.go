package codeforces

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/codeforces/list"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [SPECIFIER]",
	Short: "Lists specified data in tabular form",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if given args is a valid specifier.
		problemCnf := util.LoadLocalConf(confSettings)
		if _, err := parseSpecifier(args, problemCnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Check if mode flag is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "c" && modeFlag != "d" && modeFlag != "s" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		// Check if count is not greater than 100.
		countFlag := cmd.Flags().MustGetUint("count")
		if modeFlag == "c" && countFlag > 100 {
			return fmt.Errorf("invalid flags - flag 'count' must be in range [1, 100]")
		} else if modeFlag == "s" && countFlag > 50 {
			return fmt.Errorf("invalid flags - flag 'count' must be in range [1, 50]")
		}

		// Flag username value is defined.
		if cmd.Flags().Changed("username") &&
			modeFlag != "s" {
			// Username flag doesn't match with given mode.
			return fmt.Errorf("invalid flags - flag 'username' not applicable for mode '%v'", modeFlag)
		}

		if cmd.Flags().Changed("count") &&
			modeFlag != "c" && modeFlag != "s" {
			// Count flag doesn't match with given mode.
			return fmt.Errorf("invalid flags - flag 'count' not applicable for mode '%v'", modeFlag)
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
		countFlag := cmd.Flags().MustGetUint("count")

		arg, _ := parseSpecifier(args, problemCnf)
		list.List(arg, modeFlag, usernameFlag, countFlag)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// All flags available to command.
	listCmd.Flags().StringP("mode", "m", "c", "mode to select data to output")
	listCmd.Flags().String("username", "", "user to fetch submissions of")
	listCmd.Flags().UintP("count", "n", 5, "maximum count of rows to display")

	// All custom completions for command flags.
	listCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"c\tcontests data",
			"d\tdashboard data",
			"s\tsubmissions data",
		}
		return modes, cobra.ShellCompDirectiveDefault
	})
}
