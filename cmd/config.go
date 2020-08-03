package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kballard/go-shellquote"
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
			"Add new code template",
			"Set default template",
			"Remove template",
			"Run 'gen' after 'fetch'",
			"Set application proxy",
			"Generate tab autocompletion",
		},
	}, &idx)
	util.SurveyErr(err)

	switch idx {
	case 0:
		fmt.Println("Welcome to the template creation wizard!")
		fmt.Println("Visit cpt wiki for a comprehensive guide")

		var alias string
		err := survey.AskOne(&survey.Input{
			Message: "Template name:",
			Help: "What do you want to call this template? Ex: 'default c++'\n" +
				"Should be unique and less than 15 characters long.",
		}, &alias, survey.WithValidator(func(ans interface{}) error {
			if len(ans.(string)) == 0 {
				return fmt.Errorf("value is required")
			}

			for alias := range viper.GetStringMap("templates") {
				if alias == ans.(string) {
					return fmt.Errorf("Template with given name exists")
				}
			}
			return nil
		}))
		util.SurveyErr(err)

		newTmplt := make(map[string]interface{})
		err = survey.Ask([]*survey.Question{
			{
				Name: "file",
				Prompt: &survey.Input{
					Message: "Relative/absolute path to template file:",
				},
				Validate: func(ans interface{}) error {
					path, err := homedir.Expand(ans.(string))
					if err != nil {
						return err
					}

					if file, err := os.Stat(path); err != nil || file.IsDir() == true {
						return fmt.Errorf("Path does not correspond to a valid file")
					}
					return nil
				},
				Transform: func(ans interface{}) interface{} {
					path, _ := homedir.Expand(ans.(string))
					path, _ = filepath.Abs(path)
					return path
				},
			},
			{
				Name: "prescript",
				Prompt: &survey.Input{
					Message: "Prescript:",
				},
				Validate: func(ans interface{}) error {
					_, err := shellquote.Split(ans.(string))
					return err
				},
			},
			{
				Name: "script",
				Prompt: &survey.Input{
					Message: "Script:",
				},
				Validate: func(ans interface{}) error {
					cmdArgs, err := shellquote.Split(ans.(string))
					if len(cmdArgs) == 0 {
						return fmt.Errorf("value is required")
					}
					return err
				},
			},
			{
				Name: "postscript",
				Prompt: &survey.Input{
					Message: "Postscript:",
				},
				Validate: func(ans interface{}) error {
					_, err := shellquote.Split(ans.(string))
					return err
				},
			},
		}, &newTmplt)
		util.SurveyErr(err)

		langMap := make(map[string]string)
		err = survey.Ask([]*survey.Question{
			{
				Name: "codeforces",
				Prompt: &survey.Select{
					Message: "Language (codeforces):",
					Options: util.ExtractMapKeys(codeforces.LanguageID),
				},
			},
		}, &langMap)
		util.SurveyErr(err)

		newTmplt["languages"] = langMap
		viper.Set("templates."+alias, newTmplt)

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

	case 2:
		if len(viper.GetStringMap("templates")) == 0 {
			fmt.Println("No configured templates found")
			os.Exit(1)
		}

		data := viper.GetStringMap("templates")
		var choice string
		err := survey.AskOne(&survey.Select{
			Message: "Select template to remove",
			Options: util.ExtractMapKeys(data),
		}, &choice)
		util.SurveyErr(err)

		delete(data, choice)
		viper.Set("templates", data)

	case 3:
		var choice bool
		util.SurveyErr(survey.AskOne(&survey.Confirm{
			Message: "Run 'gen' after 'fetch'?",
			Default: false,
		}, &choice))
		viper.Set("gen_on_fetch", choice)

	case 4:
		var choice string
		err := survey.AskOne(&survey.Input{
			Message: "Proxy URL:",
			Help: "Set new proxy (should match protocol://host[:port])\n" +
				"Leave blank to reset to environment proxy",
		}, &choice, survey.WithValidator(func(ans interface{}) error {
			if ans.(string) == "" {
				return nil
			}
			_, err := url.ParseRequestURI(ans.(string))
			return err
		}))
		util.SurveyErr(err)

		prxy, _ := url.ParseRequestURI(choice)
		viper.Set("proxy_url", prxy)

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
				fmt.Println("OS", runtime.GOOS, "is not supported for bash completions")
				os.Exit(0)
			}

		case "zsh":
			err = rootCmd.GenZshCompletionFile("/usr/share/zsh/site-functions/_cpt")

		case "fish":
			gflPath, _ := homedir.Expand("~/.config/fish/completions/yourprogram.fish")
			err = rootCmd.GenFishCompletionFile(gflPath, true)

		case "powershell":
			fmt.Println("Completion script shall be written to file cpt.ps1 in current directory")
			fmt.Println("Read https://stackoverflow.com/a/20415779/9606036 for instructions to source the script")
			err = rootCmd.GenPowerShellCompletionFile("cpt.ps1")
		}

		if errors.Is(err, os.ErrPermission) {
			fmt.Println("Insufficient permissions! Try again as sudo/admin")
			fmt.Println(err)
			os.Exit(1)
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Completion scripts written successfully!")
		fmt.Println("Reload your shell for completion to take effect")
		os.Exit(0)
	}

	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Failed to save configurations")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Configurations successfully saved!")
}
