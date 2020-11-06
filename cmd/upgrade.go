package cmd

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/upgrade"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Args:  cobra.NoArgs,
	Short: "Upgrade cli to latest version",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if mode is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "s" && modeFlag != "c" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		mode := cmd.Flags().MustGetString("mode")
		switch mode {
		case "s":
			upgrade.Self(rootCmd.Version)
		case "c":
			upgrade.Checkers(checkerDir, cnf)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// All flags available to command.
	upgradeCmd.Flags().StringP("mode", "m", "s", "package to be upgraded")

	// All custom completes for command flags.
	upgradeCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"s\tself",
			"c\tcheckers",
		}

		return modes, cobra.ShellCompDirectiveDefault
	})
}

// We don't load local configurations in this command. Reasons are:
// Upgrade is a global command, with global changes.
// local settings would potentially mess up the configurations.
