package cmd

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/generate"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Create file using template",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Handle case where '--template' is not set.
		if cmd.Flag("template").Value.String() == "" {
			defaultTemplate := confSettings.GetString("generate.defaultTemplate")
			if defaultTemplate == "" {
				return fmt.Errorf("invalid flags - no template value provided")
			}
			// Set value of '--template' to defaultTemplate.
			cmd.Flag("template").Value.Set(defaultTemplate)
		}

		// Check if '--template' value is valid.
		templateFlag := cmd.Flag("template").Value.String()
		if confSettings.Get("template."+templateFlag) == nil {
			return fmt.Errorf("invalid flags - template %v not present", templateFlag)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// Local (folder) configurations to use.
		problemCnf := conf.New()
		problemCnf.LoadFile("meta.yaml")
		problemCnf.LoadDefault(confSettings.GetAll())

		templateFlag := cmd.Flag("template").Value.String()
		generate.Generate(templateFlag, problemCnf)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// All flags available to command.
	generateCmd.Flags().StringP("template", "t", "", "alias name of the template to use")

	// All custom completions for command flags.
	generateCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		aliases := confSettings.GetMapKeys("template")
		return aliases, cobra.ShellCompDirectiveDefault
	})
}
