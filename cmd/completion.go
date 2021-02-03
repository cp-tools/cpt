package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To generate shell completions:

Bash:
------------
To load completions for each session, execute once:
Linux:
	cpt completion bash > /etc/bash_completion.d/cpt
MacOS:
	cpt completion bash > /usr/local/etc/bash_completion.d/cpt

Zsh:
------------
If shell completion is not enabled in your environment, you will
need to enable it. Execute the following once:
	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for each session, execute once:
	cpt completion zsh > "${fpath[1]}/_cpt" 

Fish:
------------
To load completions for each session, execute once:
	cpt completion fish > ~/.config/fish/completions/cpt.fish

Powershell:
------------
To load completions for each session, execute once:
	cpt completion powershell > cpt.ps1

Source the generated file from your powershell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
