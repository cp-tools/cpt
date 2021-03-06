package config

import (
	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

	"github.com/AlecAivazis/survey/v2"
)

// SetStdoutColor configures user interface coloring.
func SetStdoutColor(cnf *conf.Conf) {
	choice := true
	err := survey.AskOne(&survey.Confirm{
		Message: "Do you want color printing of verbose text?",
		Default: true,
	}, &choice)
	utils.SurveyOnInterrupt(err)

	cnf.Set("ui.stdoutColor", choice)
}
