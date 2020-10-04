package cmd

import (
	"github.com/cp-tools/cpt/cmd/internal/config"

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
			Message: "Which configuration do you want to perform?",
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
			config.AddTemplate(confTemplates)
			confTemplates.WriteFile()
		case 1:
			config.RemoveTemplate(confTemplates)
			confTemplates.WriteFile()
		case 2:
			config.SetGenerateOnFetch(confSettings)
			confSettings.WriteFile()
		case 3:
			aliases := confTemplates.GetMapKeys("")
			config.SetDefaultTemplate(confSettings, aliases)
			confSettings.WriteFile()
		case 4:
			config.SetHeadlessBrowser(confSettings)
			confSettings.WriteFile()
		case 5:
			config.SetStdoutColor(confSettings)
			confSettings.WriteFile()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
