package cmd

import (
	"github.com/cp-tools/cpt/cmd/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.NoArgs,
	Short: "Configure global settings",
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt user for configuration to modify.
		index := 0
		survey.AskOne(&survey.Select{
			Message: "What configuration do you want to perform?",
			Options: []string{
				"template - add new",
				"template - delete existing",
				"generate - run on 'fetch'",
				"generate - set default template",
				"browser - set headless browser",
				"ui - set stdout colorization",
			},
		}, &index)

		switch index {
		case 0:
			config.AddTemplate(confSettings)
		case 1:
			config.RemoveTemplate(confSettings)
		case 2:
			config.SetGenerateOnFetch(confSettings)
		case 3:
			config.SetDefaultTemplate(confSettings)
		case 4:
			config.SetHeadlessBrowser(confSettings)
		case 5:
			config.SetStdoutColor(confSettings)
		}
		// Write file after changes are done.
		confSettings.WriteFile()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
