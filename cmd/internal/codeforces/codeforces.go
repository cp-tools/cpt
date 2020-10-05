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

// InitConfSettings merges configurations from global settings.
func InitConfSettings(confDir string, confMap map[string]interface{}) {
	// Load global configurations.
	confSettings.Load(confMap)
	// Load local configurations (overwrite global values).
	confSettingsPath := filepath.Join(confDir, "codeforces.yaml")
	confSettings.LoadFile(confSettingsPath)
}
