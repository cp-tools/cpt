package codeforces

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt/cmd/codeforces/submit"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [SPECIFIER]",
	Short: "Submit problem solution to judge",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if given args is a valid specifier.
		problemCnf := util.LoadLocalConf(confSettings)
		if _, err := parseSpecifier(args, problemCnf); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}

		// Check if given file path points to valid file.
		fileFlag := cmd.Flags().MustGetString("file")
		if fileFlag != "" {
			if file, err := os.Stat(fileFlag); os.IsNotExist(err) || file.IsDir() {
				return fmt.Errorf("invalid flags - %v is not a valid file", fileFlag)
			}
		}

		// Load headless browser to use.
		startHeadlessBrowser()

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// Local (folder) configurations to use.
		problemCnf := util.LoadLocalConf(confSettings)

		fileFlag := cmd.Flag("file").Value.String()

		arg, _ := parseSpecifier(args, problemCnf)
		submit.Submit(arg, fileFlag, problemCnf)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// All flags available to command.
	submitCmd.Flags().StringP("file", "f", "", "code file to submit")
}
