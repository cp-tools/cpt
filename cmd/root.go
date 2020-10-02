package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

var rootCmd = &cobra.Command{
	Use:     "cpt",
	Short:   "Lightweight cli tool for competitive programming!",
	Version: "v0.12.1",
}

var (
	configDir       string
	configSettings  *koanf.Koanf
	configTemplates *koanf.Koanf
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
		initConfigDir,
		initSettings,
		initTemplates,
	)

	// set OnSIGINT function for survey
	survey.OnSIGINTFunc = func() {
		fmt.Println("interrupted")
		os.Exit(1)
	}
}

// determine and set configDir path
func initConfigDir() {
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir = filepath.Join(dir, "cp-tools", "cpt")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		log.Fatalf("error creating config folder: %v", err)
	}
}

func initSettings() {
	configSettings = koanf.New(".")
	// configure default values
	configSettings.Load(confmap.Provider(map[string]interface{}{
		"ui.stdoutColor": true,
	}, "."), nil)

	configSettingsPath := filepath.Join(configDir, "cpt.yaml")
	if _, err := os.Stat(configSettingsPath); os.IsNotExist(err) {
		if _, err := os.Create(configSettingsPath); err != nil {
			log.Fatalf("error creating settings file: %v", err)
		}
	}

	// load YAML settings config.
	if err := configSettings.Load(file.Provider(configSettingsPath), yaml.Parser()); err != nil {
		log.Fatalf("error loading settings file: %v", err)
	}
}

func initTemplates() {
	configTemplates = koanf.New(".")

	configTemplatesPath := filepath.Join(configDir, "templates.yaml")
	if _, err := os.Stat(configTemplatesPath); os.IsNotExist(err) {
		if _, err := os.Create(configTemplatesPath); err != nil {
			log.Fatalf("error creating templates file: %v", err)
		}
	}

	// load YAML templates config.
	if err := configTemplates.Load(file.Provider(configTemplatesPath), yaml.Parser()); err != nil {
		log.Fatalf("error loading templates file: %v", err)
	}
}
