package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/packages/conf"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cpt",
	Short:   "Lightweight cli tool for competitive programming!",
	Version: "v0.12.1",
}

var (
	confDir       string
	confSettings  *conf.Conf
	confTemplates *conf.Conf
)

// Execute adds all child commands to the root command and
// sets flags appropriately. Called by main.main()
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		initConfDir,
		initConfSettings,
		initConfTemplates,
	)

	// Set OnSIGINT function for survey module.
	survey.OnSIGINTFunc = func() {
		fmt.Println("interrupted")
		os.Exit(1)
	}
}

// Determine and set configDir path.
func initConfDir() {
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	confDir = filepath.Join(dir, "cp-tools", "cpt")
	if err := os.MkdirAll(confDir, os.ModePerm); err != nil {
		log.Fatalf("error creating config folder: %v", err)
	}
}

func initConfSettings() {
	confSettings = conf.New()
	// Configure default values.
	confSettings.Set("ui.stdoutColor", true)

	confSettingsPath := filepath.Join(confDir, "cpt.yaml")
	confSettings.LoadFile(confSettingsPath)
}

func initConfTemplates() {
	confTemplates = conf.New()

	confTemplatesPath := filepath.Join(confDir, "templates.yaml")
	confTemplates.LoadFile(confTemplatesPath)
}
