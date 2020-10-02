package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/cp-tools/cpt/util"
	"github.com/fatih/color"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cptConfigCmd represents the config command
var cptConfigCmd = &cobra.Command{
	Use:  "config",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cptConfig()
	},
	Short: "Configure application wide settings",
	Long: `Application wide settings configuration.
Options like 'colorization', 'template generation', 'proxy' etc
are configured here. Use the selection menu that appears to make 
changes to the settings.

All settings are saved to $CONFIGDIR/cpt/cpt.json file.
`,
}

func init() {
	rootCmd.AddCommand(cptConfigCmd)
}

func cptConfig() {
	var idx int
	err := survey.AskOne(&survey.Select{
		Message: "Select configuration:",
		Options: []string{
			"Set default template",
			"Run 'gen' after 'fetch'",
			"Set application proxy",
			"Generate tab autocompletion",
			"Set interface colorization",
		},
	}, &idx)
	util.SurveyErr(err)

	switch idx {

	case 1:
		opts := append(util.ExtractMapKeys(viper.GetStringMap("templates")), "none")
		var choice string
		err := survey.AskOne(&survey.Select{
			Message: "Select template:",
			Default: viper.GetString("default_template"),
			Options: opts,
		}, &choice)
		util.SurveyErr(err)

		viper.Set("default_template", choice)

	case 3:
		var choice bool
		util.SurveyErr(survey.AskOne(&survey.Confirm{
			Message: "Run 'gen' after 'fetch'?",
			Default: false,
		}, &choice))
		viper.Set("gen_on_fetch", choice)

	case 5:
		var choice string
		util.SurveyErr(survey.AskOne(&survey.Select{
			Message: "Select shell type:",
			Options: []string{"bash", "zsh", "fish", "powershell"},
		}, &choice))

		var err error
		switch choice {
		case "bash":
			switch runtime.GOOS {
			case "linux":
				err = rootCmd.GenBashCompletionFile("/etc/bash_completion.d/cpt")
			case "darwin":
				err = rootCmd.GenBashCompletionFile("/usr/local/etc/bash_completion.d/cpt")
			default:
				color.Yellow("OS %v is not supported for bash completions", runtime.GOOS)
				os.Exit(0)
			}

		case "zsh":
			err = rootCmd.GenZshCompletionFile("/usr/share/zsh/site-functions/_cpt")

		case "fish":
			gflPath, _ := homedir.Expand("~/.config/fish/completions/yourprogram.fish")
			err = rootCmd.GenFishCompletionFile(gflPath, true)

		case "powershell":
			color.Blue("Completion script shall be written to file cpt.ps1 in current directory")
			color.Blue("Read https://stackoverflow.com/a/20415779/9606036 for instructions to source the script")
			err = rootCmd.GenPowerShellCompletionFile("cpt.ps1")
		}

		if errors.Is(err, os.ErrPermission) {
			color.Red("Insufficient permissions! Try again as sudo/admin")
			fmt.Println(err)
			os.Exit(1)
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		color.Green("Completion scripts written successfully!")
		color.Green("Reload your shell for completion to take effect")
		os.Exit(0)

	}

	if err := viper.WriteConfig(); err != nil {
		color.Red("Failed to save configurations")
		fmt.Println(err)
		os.Exit(1)
	}
	color.Green("Configurations successfully saved!")
}
