package codeforces

import (
	"path/filepath"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var (
	confSettings = conf.New()
)

// SetParentCmd sets parent command of all subcommands
// in this module to parentCmd.
func SetParentCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rootCmd.Commands()...)
	// Set rootCmd to parentCmd.
	rootCmd = parentCmd
}

// ConfLoadFile loads codeforces.yaml from specified directory.
func ConfLoadFile(confDir string) {
	confSettingsPath := filepath.Join(confDir, "codeforces.yaml")
	confSettings.LoadFile(confSettingsPath)
}

// ConfLoadDefaults sets default values in local module.
func ConfLoadDefaults(confMap map[string]interface{}) {
	confSettings.LoadDefault(confMap)
	// Set local defaults here.

	// Path structure when 'fetching' problem tests.
	confSettings.SetDefault("fetch.problemFolderPath", []string{
		"codeforces", "{{.Arg.Contest}}", "{{.Arg.Problem}}",
	})
	// Path structure when 'pulling' problem submissions.
	confSettings.SetDefault("pull.problemFolderPath", []string{
		"codeforces", "{{.Arg.Contest}}", "{{.Arg.Problem}}",
	})

}

func startHeadlessBrowser() {
	binary := confSettings.GetString("browser.binary")
	profile := confSettings.GetString("browser.profile")
	codeforces.Start(true, profile, binary)
}

func parseSpecifier(args []string, cnf *conf.Conf) (codeforces.Args, error) {
	arg, err := codeforces.Parse(strings.Join(args, ""))
	if err != nil {
		return arg, err
	}

	if arg == (codeforces.Args{}) && cnf.Has("problem.arg") {
		// Parse from configuration.
		err = mapstructure.Decode(cnf.Get("problem.arg"), &arg)
	}
	return arg, err
}
