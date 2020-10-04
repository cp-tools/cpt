package config

import (
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/AlecAivazis/survey/v2"
)

// SetGenerateOnFetch sets 'generate.onFetch'.
func SetGenerateOnFetch(cnf *conf.Conf) {
	choice := false
	survey.AskOne(&survey.Confirm{
		Message: "Do you want to generate the default template when problem tests are fetched?",
		Default: false,
	}, &choice)

	cnf.Set("generate.onFetch", choice)
}

// SetDefaultTemplate sets 'generate.defaultTemplate'.
func SetDefaultTemplate(cnf *conf.Conf, aliases []string) {
	alias := ""
	survey.AskOne(&survey.Select{
		Message: "Which template do you want to make the default?",
		Options: append(aliases, ""),
	}, &alias)

	if alias == "" {
		// Remove default template value.
		cnf.Erase("generate.defaultTemplate")
		return
	}

	cnf.Set("generate.defaultTemplate", alias)
}
