package cmd

import (
	"github.com/cp-tools/cpt/cmd/codeforces"
	"github.com/spf13/cobra"
)

var codeforcesCmd = &cobra.Command{
	Use:     "codeforces",
	Aliases: []string{"cf"},
	Short:   "Functions exclusive to codeforces",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load and initialize conf related settings.
		codeforces.InitModuleConf(cnf, confDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(codeforcesCmd)
	// Load subcommands into codeforcesCmd.
	codeforces.SetParentCmd(codeforcesCmd)
}
