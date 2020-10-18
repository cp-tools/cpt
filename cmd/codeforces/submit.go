package codeforces

import (
	"fmt"
	"os"
	"strings"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/cmd/codeforces/submit"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [SPECIFIER]",
	Short: "Submit problem solution to judge",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load headless browser to use.
		startHeadlessBrowser()

		// Check if given args is a valid specifier.
		if _, err := codeforces.Parse(strings.Join(args, "")); err != nil {
			return fmt.Errorf("invalid args - %v", err)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
		// Local (folder) configurations to use.
		problemCnf := conf.New()
		problemCnf.LoadFile("meta.yaml")
		problemCnf.LoadDefault(confSettings.GetAll())

		fileFlag := cmd.Flag("file").Value.String()
		// Parse args to specifier.
		arg, _ := codeforces.Parse(strings.Join(args, ""))
		submit.Submit(arg, fileFlag, problemCnf)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// All flags available to command.
	submitCmd.Flags().StringP("file", "f", "", "code file to submit")
}
