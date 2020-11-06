package codeforces

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var (
	cnf *conf.Conf
)

// SetParentCmd sets parent command of all subcommands
// in this module to parentCmd.
func SetParentCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rootCmd.Commands()...)
	// Set rootCmd to parentCmd.
	rootCmd = parentCmd
}

// InitModuleConf sets codeforces configurations.
func InitModuleConf(rootCnf *conf.Conf, confDir string) {
	// Set default values here.
	cnf = conf.New("codeforces").SetParent(rootCnf)
	cnf.SetDefault("fetch.problemFolderPath", []string{
		"codeforces", "{{.Arg.Contest}}", "{{.Arg.Problem}}",
	})
	cnf.SetDefault("pull.problemFolderPath", []string{
		"codeforces", "{{.Arg.Contest}}", "{{.Arg.Problem}}",
	})

	cnfFilePath := filepath.Join(confDir, "codeforces.yaml")
	cnf.LoadFile(cnfFilePath)
}

func startHeadlessBrowser() {
	binary := cnf.GetString("browser.binary")
	profile := cnf.GetString("browser.profile")
	if err := codeforces.Start(true, profile, binary); err != nil {
		fmt.Println(color.RedString("error launching automated browser:"), err)
		os.Exit(1)
	}
}

func parseSpecifier(args []string, rootCnf *conf.Conf) (codeforces.Args, error) {
	arg, err := codeforces.Parse(strings.Join(args, ""))
	if err != nil {
		return arg, err
	}

	if arg == (codeforces.Args{}) && rootCnf.Has("problem.arg") {
		// Parse from configuration.
		err = mapstructure.Decode(rootCnf.Get("problem.arg"), &arg)
	}
	return arg, err
}
