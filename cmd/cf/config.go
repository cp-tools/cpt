package cf

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"
	"github.com/infixint943/cookiejar"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cfConfigCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings exclusive to codeforces",
	Run: func(cmd *cobra.Command, args []string) {
		config()
	},
	Args: cobra.NoArgs,
}

func init() {
	RootCmd.AddCommand(configCmd)
}

func config() {
	var choice int
	err := survey.AskOne(&survey.Select{
		Message: "Select configuration:",
		Options: []string{
			"Login to codeforces",
		},
	}, &choice)
	util.SurveyErr(err)

	switch choice {
	case 0:
		// check if any saved session present
		if usr := cfViper.GetString("username"); len(usr) != 0 {
			fmt.Println(color.BlueString("Current user handle:"), usr)
			color.Yellow("Existing session will be OVERWRITTEN!")
		}

		var usr, passwd string
		err := survey.AskOne(&survey.Input{Message: "Username:"},
			&usr, survey.WithValidator(survey.Required))
		util.SurveyErr(err)
		err = survey.AskOne(&survey.Password{Message: "Password:"},
			&passwd, survey.WithValidator(survey.Required))
		util.SurveyErr(err)

		color.Blue("Logging in. Please wait.")
		// remove all past session cookies
		jar, _ := cookiejar.New(nil)
		codeforces.SessCln.Jar = jar

		handle, err := codeforces.Login(usr, passwd)
		if err != nil {
			color.Red("Login failed")
			fmt.Println(err)
			os.Exit(0)
		}

		color.Green("Login successful")
		fmt.Println(color.BlueString("Welcome"), handle)

		cfViper.Set("username", usr)
		if ed, err := util.Encrypt(usr, passwd); err == nil {
			cfViper.Set("password", ed)
		} else {
			color.Red("Could not encrypt password")
			fmt.Println(err)
			os.Exit(1)
		}

		hostURL, _ := url.Parse("https://codeforces.com")
		http.DefaultClient.Jar.SetCookies(hostURL, jar.Cookies(hostURL))
		viper.Set("cookies", http.DefaultClient.Jar)

		cfViper.WriteConfig()
		viper.WriteConfig()
	}
}
