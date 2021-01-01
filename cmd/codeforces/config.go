package codeforces

import (
	"github.com/cp-tools/cpt-lib/v2/codeforces"
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
		err := survey.AskOne(&survey.Select{
			Message: "What configuration do you want to perform?",
			Options: []string{
				"template - set language",
				"generate - run on 'fetch'",
				"generate - set default template",
				"browser - set headless browser",
			},
		}, &index)
		util.SurveyOnInterrupt(err)

		rootCnf := cnf.GetParent("codeforces")
		switch index {
		case 0:
			languages := util.ExtractMapKeys(codeforces.LanguageID)
			config.SetTemplateLanguage(rootCnf, languages)

		case 1:
			config.SetGenerateOnFetch(rootCnf)

		case 2:
			config.SetDefaultTemplate(rootCnf)

		case 3:
			config.SetHeadlessBrowser(rootCnf)
		}
		// Write file after changes are done.
		rootCnf.WriteFile()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
