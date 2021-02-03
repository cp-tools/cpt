package cmd

import (
	"fmt"

	"github.com/cp-tools/cpt/cmd/generate"
	"github.com/cp-tools/cpt/pkg/conf"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Create file using template",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cnf = conf.New("local").SetParent(cnf).LoadFile("meta.yaml")

		// Handle case where '--template' is not set.
		if cmd.Flags().MustGetString("template") == "" {
			defaultTemplate := cnf.GetString("generate.defaultTemplate")
			if defaultTemplate == "" {
				return fmt.Errorf("invalid flags - no template value provided")
			}
			// Set value of '--template' to defaultTemplate.
			cmd.Flags().Set("template", defaultTemplate)
		}

		// Check if '--template' value is valid.
		templateFlag := cmd.Flags().MustGetString("template")
		if cnf.Get("template."+templateFlag) == nil {
			return fmt.Errorf("invalid flags - template %v not present", templateFlag)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		template := cmd.Flags().MustGetString("template")
		generate.Generate(template, cnf, cnf.GetAll())
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// All flags available to command.
	generateCmd.Flags().StringP("template", "t", "", "alias name of the template to use")

	// All custom completions for command flags.
	generateCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		aliases := cnf.GetMapKeys("template")
		return aliases, cobra.ShellCompDirectiveDefault
	})
}
