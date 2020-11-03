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
		cnf = util.LoadLocalConf(cnf)

		// Check if given args is a valid specifier.
		if _, err := parseSpecifier(args, cnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Check if mode flag is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "c" && modeFlag != "d" && modeFlag != "s" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		// Flag username value is defined.
		if cmd.Flags().Changed("username") && modeFlag != "s" {
			// Username flag doesn't match with given mode.
			return fmt.Errorf("invalid flags - flag 'username' not applicable for mode '%v'", modeFlag)
		}

		if cmd.Flags().Changed("count") && modeFlag != "c" && modeFlag != "s" {
			// Count flag doesn't match with given mode.
			return fmt.Errorf("invalid flags - flag 'count' not applicable for mode '%v'", modeFlag)
		}

		// Load headless browser to use.
		startHeadlessBrowser()

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		username := cmd.Flags().MustGetString("username")
		mode := cmd.Flags().MustGetString("mode")
		count := cmd.Flags().MustGetUint("count")

		arg, _ := parseSpecifier(args, cnf)

		switch mode {
		case "c":
			list.Contests(arg, count)
		case "d":
			list.Dashboard(arg)
		case "s":
			list.Submissions(arg, username, count)
		}
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
			"c\tcontests",
			"d\tdashboard",
			"s\tsubmissions",
		}
		return modes, cobra.ShellCompDirectiveDefault
	})
}
