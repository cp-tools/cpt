package codeforces

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list [SPECIFIER]",
	Short: "Lists specified data in tabular form",
}

func init() {
	rootCmd.AddCommand(listCmd)

	// All flags available to command.
	listCmd.Flags().StringP("mode", "m", "la", "mode to select data to output")
	listCmd.Flags().String("username", "", "user to fetch submissions of")

	// All custom completions for command flags.
	listCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		modes := []string{
			"c\tcontests data",
			"d\tdashboard data",
			"s\tsubmissions data",
		}
		return modes, cobra.ShellCompDirectiveDefault
	})
}
