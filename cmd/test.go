package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cp-tools/cpt/cmd/test"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run code file against sample tests",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		util.LoadLocalConf(cnf)

		// Check if mode is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "j" && modeFlag != "i" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		// Check if checker is valid.
		checkerFlag := cmd.Flags().MustGetString("checker")
		// First check if checker exists in relative directory.
		if file, err := os.Stat(checkerFlag); os.IsNotExist(err) || file.IsDir() {
			// Checker doesn't exist in relative path. Search in cpt-checker folder.
			checkerPath := filepath.Join(rootDir, "cpt-checker", checkerFlag)
			if file, err := os.Stat(checkerPath); os.IsNotExist(err) || file.IsDir() {
				// Checker specified is invalid.
				return fmt.Errorf("invalid flags - checker '%v' doesn't exist", checkerFlag)
			}
			// Update checker variable.
			cmd.Flags().Set("checker", checkerPath)
		}

		// Check if given file path point to valid file.
		fileFlag := cmd.Flags().MustGetString("file")
		if fileFlag != "" {
			if file, err := os.Stat(fileFlag); os.IsNotExist(err) || file.IsDir() {
				return fmt.Errorf("invalid flags - %v is not a valid file", fileFlag)
			}
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		checker := cmd.Flags().MustGetString("checker")
		filePath := cmd.Flags().MustGetString("file")
		mode := cmd.Flags().MustGetString("mode")
		timelimit := cmd.Flags().MustGetDuration("timelimit")

		test.Test(checker, filePath, mode, timelimit, cnf)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// All flags available to command.
	testCmd.Flags().StringP("checker", "c", "lcmp", "testlib checker to use")
	testCmd.Flags().StringP("file", "f", "", "code file to run tests on")
	testCmd.Flags().StringP("mode", "m", "j", "mode to run tests on")
	testCmd.Flags().DurationP("timelimit", "t", 2*time.Second, "timelimit per test")

	// All custom completions for command flags.
	testCmd.RegisterFlagCompletionFunc("checker", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		checkers := cnf.GetMapKeys("checker")
		for i := range checkers {
			desc := cnf.GetString("checker." + checkers[i] + ".desc")
			checkers[i] = fmt.Sprintf("%v\t%v", checkers[i], desc)
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
