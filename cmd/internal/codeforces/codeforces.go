package codeforces

import (
	"path/filepath"

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

// InitConfSettings loads configurations from
// global and local settings.
func InitConfSettings(confDir string, cnfMap map[string]interface{}) {
	// Load global configurations into confSettings.
	confSettings.Load(cnfMap)
	// Load local configurations (overwrite global values).
	confSettingsPath := filepath.Join(confDir, "codeforces.yaml")
	confSettings.LoadFile(confSettingsPath)
}
