package config

import "github.com/AlecAivazis/survey/v2"

func setStdoutColor(dataMap map[string]interface{}) {
	choice := true
	survey.AskOne(&survey.Confirm{
		Message: "Do you want color printing of verbose text?",
		Default: true,
	}, &choice)
	dataMap["ux"] = map[string]bool{"stdoutColor": choice}
	return
}

func setHeadlessBrowser(dataMap map[string]interface{}) {
}
