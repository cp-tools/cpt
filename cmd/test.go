package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cp-tools/cpt/packages/conf"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run code file against sample tests",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// All flags available to command.
	testCmd.Flags().StringP("checker", "c", "lcmp", "testlib checker to use")
	testCmd.Flags().StringP("file", "f", "", "code file to run tests on")
	testCmd.Flags().StringP("mode", "m", "j", "mode to run tests on")
	testCmd.Flags().DurationP("time-limit", "t", 2*time.Second, "time limit per test")

	// All custom completions for command flags.
	testCmd.RegisterFlagCompletionFunc("checker", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Find all executables in cpt-checkers directory.
		confCheckerPath := filepath.Join(rootDir, "cpt-checker", "checkers.yaml")
		confChecker := conf.New()
		confChecker.LoadFile(confCheckerPath)
		// Add description to tab completion.
		checkers := confChecker.GetMapKeys("")
		for i := 0; i < len(checkers); i++ {
			desc := confChecker.GetString(checkers[0] + ".description")
			checkers[0] += fmt.Sprintf("\t%v", desc)
		}
		return checkers, cobra.ShellCompDirectiveDefault
	})

	testCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"j\tjudge mode",
			"i\tinteractive mode",
		}

		return modes, cobra.ShellCompDirectiveDefault
	})
}
