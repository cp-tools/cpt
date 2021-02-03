package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/packages/conf"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cpt",
	Short: "Lightweight cli tool for competitive programming!",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize global configurations.
		initGlobalConf()
		// Set verbose text coloring.
		color.NoColor = !cnf.GetBool("ui.stdoutColor")
	},

	Version: "v0.13.2",

	TraverseChildrenHooks: true,
}

var (
	rootDir    string
	confDir    string
	checkerDir string

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

	checkerDir = filepath.Join(rootDir, "cpt-checker")
	if err := os.MkdirAll(checkerDir, os.ModePerm); err != nil {
		log.Fatalf("error creating checker folder: %v", err)
	}
}

func initGlobalConf() {
	// Load global configuration.
	cnf = conf.New("global")
	cnf.SetDefault("ui.stdoutColor", true)

	cnfFilePath := filepath.Join(confDir, "cpt.yaml")
	cnf.LoadFile(cnfFilePath)

	// Load checker configuration.
	cnf = conf.New("checker").SetParent(cnf)

	cnfFilePath = filepath.Join(checkerDir, "checker.yaml")
	cnf.LoadFile(cnfFilePath)
}

/*
First comes global conf.
Next comes all global, non main confs.
Then sub module (codeforces, atcoder) conf.
Lastly, local (folder) conf.

global --> checker -->  .... --> local
*/
