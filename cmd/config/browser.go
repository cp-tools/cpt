package config

import (
	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
)

// SetHeadlessBrowser configures headless browser to use.
func SetHeadlessBrowser(cnf *conf.Conf) {
	browserMap := make(map[string]interface{})
	err := survey.Ask([]*survey.Question{
		{
			Name: "binary",
			Prompt: &survey.Input{
				Message: "What is the command to start a new browser instance?",
				Help: `
For website related functionalities, the browser specified
above is run (in headless mode) to perform various tasks.
Use the same browser you use to access the corresponding sites,
as cookies and login details are shared by this tool.

Currently, only google chrome/chromium and edge are supported.

Example commands (without quotes) to run supported browsers are:
- 'google-chrome' (linux)
- '/usr/bin/google-chrome' (linux)
- 'C:\Program Files (x86)\Google\Chrome\Application\chrome.exe' (windows)
- 'C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe' (windows)
- '/Applications/Chromium.app/Contents/MacOS/Chromium' (darwin)
- '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' (darwin)
`,
			},
			Validate: survey.Required,
		},
		{
			Name: "profile",
			Prompt: &survey.Input{
				Message: "Where is the default profile of the browser located?",
				Help: `
The browser profile contains all session cookies, which is required
by this tool, to access the current logged in session.

Here are the locations of browser profiles of various browsers:
- Google Chrome:
  => '/home/<username>/.config/google-chrome/' (linux)
  => 'C:\Users\<username>\AppData\Local\Google\Chrome\User Data\' (windows)
  => 'Users/<username>/Library/Application Support/Google/Chrome/' (darwin)

- Microsoft Edge
  => 'C:\Users\<username>\AppData\Local\Microsoft\Edge\User Data\' (windows)
`,
			},
			Validate: func(ans interface{}) error {
				_, err := homedir.Abs(ans.(string))
				if err != nil {
					return err
				}

				return nil
			},
			Transform: func(ans interface{}) (newAns interface{}) {
				profilePath, _ := homedir.Abs(ans.(string))
				return profilePath
			},
		},
	}, &browserMap)
	util.SurveyOnInterrupt(err)

	cnf.Set("browser", browserMap)
}
