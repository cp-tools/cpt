package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:
Bash:

# To load completions for each session, execute once:
Linux:
  $ cpt completion bash > /etc/bash_completion.d/cpt
MacOS:
  $ cpt completion bash > /usr/local/etc/bash_completion.d/cpt

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ cpt completion zsh > "${fpath[1]}/_cpt"

# You will need to start a new shell for this setup to take effect.

Fish:

$ cpt completion fish | source

# To load completions for each session, execute once:
$ cpt completion fish > ~/.config/fish/completions/cpt.fish
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
