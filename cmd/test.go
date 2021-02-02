package cmd

import (
	"fmt"
	"time"

	"github.com/cp-tools/cpt/cmd/test"
	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run code file against sample tests",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cnf = conf.New("local").SetParent(cnf).LoadFile("meta.yaml")

		// Check if mode is valid.
		modeFlag := cmd.Flags().MustGetString("mode")
		if modeFlag != "d" && modeFlag != "i" {
			return fmt.Errorf("invalid flags - unknown mode '%v'", modeFlag)
		}

		// Check if checker is valid.
		checkerFlag := cmd.Flags().MustGetString("checker")
		if !cnf.Has("checker.checkers." + checkerFlag) {
			return fmt.Errorf("invalid flags - checker '%v' not configured", checkerFlag)
		}

		// Check if given file path points to valid file.
		fileFlag := cmd.Flags().MustGetString("file")
		if fileFlag != "" {
			if !utils.FileExists(fileFlag) {
				return fmt.Errorf("invalid flags - %v is not a valid file", fileFlag)
			}
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		checker := cmd.Flags().MustGetString("checker")
		file := cmd.Flags().MustGetString("file")
		mode := cmd.Flags().MustGetString("mode")
		timeLimit := cmd.Flags().MustGetDuration("time-limit")
		memoryLimit := cmd.Flags().MustGetUint64("memory-limit")

		// If user has not specified time limit for
		// interactive testing, set 1 hour limit.
		if mode == "i" && !cmd.Flags().Changed("time-limit") {
			timeLimit = time.Hour
		}

		checkerScript := cnf.GetString("checker.checkers." + checker + ".script")

		test.Test(file, checkerScript,
			timeLimit, memoryLimit,
			"", "", mode, cnf)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// All flags available to command.
	testCmd.Flags().StringP("checker", "c", "lcmp", "testlib checker to use")
	testCmd.Flags().StringP("file", "f", "", "code file to run tests on")
	testCmd.Flags().StringP("mode", "m", "d", "mode to run tests on")
	testCmd.Flags().DurationP("time-limit", "t", 2*time.Second, "time limit per test")
	testCmd.Flags().Uint64("memory-limit", 256*1024*1024, "memory limit per test (in bytes)")

	// All custom completions for command flags.
	testCmd.RegisterFlagCompletionFunc("checker", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		checkers := cnf.GetMapKeys("checker.checkers")
		for i := range checkers {
			desc := cnf.GetString("checker.checkers." + checkers[i] + ".desc")
			checkers[i] = fmt.Sprintf("%v\t%v", checkers[i], desc)
		}

		return checkers, cobra.ShellCompDirectiveDefault
	})

	testCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"d\tdefault",
			"i\tinteractive",
		}

		return modes, cobra.ShellCompDirectiveDefault
	})
}
