package cmd

import (
	"github.com/cp-tools/cpt/cmd/config"
	"github.com/cp-tools/cpt/utils"

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
		err := survey.AskOne(&survey.Select{
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
		utils.SurveyOnInterrupt(err)

		// We are editing global config here.
		rootCnf := cnf.GetParent("global")
		switch index {
		case 0:
			config.AddTemplate(rootCnf)
		case 1:
			config.RemoveTemplate(rootCnf)
		case 2:
			config.SetGenerateOnFetch(rootCnf)
		case 3:
			config.SetDefaultTemplate(rootCnf)
		case 4:
			config.SetHeadlessBrowser(rootCnf)
		case 5:
			config.SetStdoutColor(rootCnf)
		}
		// Write file after changes are done.
		rootCnf.WriteFile()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
