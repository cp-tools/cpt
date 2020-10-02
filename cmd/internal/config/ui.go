package config

import "github.com/AlecAivazis/survey/v2"

// SetStdoutColor configures user interface coloring.
func SetStdoutColor(dataMap map[string]interface{}) {
	choice := true
	survey.AskOne(&survey.Confirm{
		Message: "Do you want color printing of verbose text?",
		Default: true,
	}, &choice)
	dataMap["ui"] = map[string]bool{"stdoutColor": choice}
	return
}
