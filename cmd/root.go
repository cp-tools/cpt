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
	Use:   "cpt",
	Short: "Lightweight cli tool for competitive programming!",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize global configurations.
		initGlobalConf()
	},

	Version: "v0.12.1",
}

var (
	rootDir string
	confDir string

	cnf *conf.Conf
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
	cobra.OnInitialize(initConfDir)
	// Set Persistent commands to be queued.
	cobra.EnablePersistentRunOverride = false

	// Set OnSIGINT function for survey module.
	survey.OnInterrupt = func() {
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
	rootDir = filepath.Join(dir, "cp-tools")
	confDir = filepath.Join(rootDir, "cpt")
	if err := os.MkdirAll(confDir, os.ModePerm); err != nil {
		log.Fatalf("error creating config folder: %v", err)
	}
}

func initGlobalConf() {
	// Load global configuration.
	cnf = conf.New("global")
	cnf.Set("ui.stdoutColor", true)

	cnfFilePath := filepath.Join(confDir, "cpt.yaml")
	cnf.LoadFile(cnfFilePath)

	// Load checker configuration.
	cnf = conf.New("checker").SetParent(cnf)

	cnfFilePath = filepath.Join(rootDir, "cpt-checker", "checker.yaml")
	cnf.LoadFile(cnfFilePath)
}

/*
First comes global conf.
Next comes all global, non main confs.
Then sub module (codeforces, atcoder) conf.
Lastly, local (folder) conf.

global --> checker -->  .... --> local
*/
