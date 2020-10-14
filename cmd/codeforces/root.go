package codeforces

import (
	"path/filepath"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/packages/conf"

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
}

func startHeadlessBrowser() {
	binary := confSettings.GetString("browser.binary")
	profile := confSettings.GetString("browser.profile")
	codeforces.Start(true, profile, binary)
	// todo: Add flags to disable image rendering.
}
