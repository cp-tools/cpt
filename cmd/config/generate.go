package config

import (
	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
)

// SetGenerateOnFetch sets 'generate.onFetch'.
func SetGenerateOnFetch(cnf *conf.Conf) {
	choice := false
	err := survey.AskOne(&survey.Confirm{
		Message: "Do you want to generate the default template when problem tests are fetched?",
		Default: false,
	}, &choice)
	util.SurveyOnInterrupt(err)

	cnf.Set("generate.onFetch", choice)
}

// SetDefaultTemplate sets 'generate.defaultTemplate'.
func SetDefaultTemplate(cnf *conf.Conf) {
	alias := ""
	err := survey.AskOne(&survey.Select{
		Message: "Which template do you want to make the default?",
		Options: append(cnf.GetMapKeys("template"), ""),
	}, &alias)
	util.SurveyOnInterrupt(err)

	if alias == "" {
		// Remove default template value.
		cnf.Delete("generate.defaultTemplate")
		return
	}

	cnf.Set("generate.defaultTemplate", alias)
}
