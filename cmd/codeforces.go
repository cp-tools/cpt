package cmd

import (
	"github.com/cp-tools/cpt/cmd/internal/codeforces"
	"github.com/spf13/cobra"
)

var codeforcesCmd = &cobra.Command{
	Use:     "codeforces",
	Aliases: []string{"cf"},
	Short:   "Functions exclusive to codeforces",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configurations into codeforces local conf.
		codeforces.InitConfSettings(confDir, confSettings.GetAll())
	},
}

func init() {
	rootCmd.AddCommand(codeforcesCmd)
	// Load subcommands into codeforcesCmd.
	codeforces.SetParentCmd(codeforcesCmd)
}
