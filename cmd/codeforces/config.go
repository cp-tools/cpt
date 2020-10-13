package codeforces

import (
	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/cmd/config"
	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.NoArgs,
	Short: "Configure codeforces settings",
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt user for configuration to modify.
		index := 0
		survey.AskOne(&survey.Select{
			Message: "What configuration do you want to perform?",
			Options: []string{
				"template - set language",
				"generate - run on 'fetch'",
				"generate - set default template",
				"browser - set headless browser",
			},
		}, &index)

		switch index {
		case 0:
			languages := util.ExtractMapKeys(codeforces.LanguageID)
			config.SetTemplateLanguage(confSettings, languages)

		case 1:
			config.SetGenerateOnFetch(confSettings)

		case 2:
			config.SetDefaultTemplate(confSettings)

		case 3:
			config.SetHeadlessBrowser(confSettings)
		}
		// Write file after changes are done.
		confSettings.WriteFile()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
