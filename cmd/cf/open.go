package cf

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [SPECIFIER]",
	Short: "Open specified problem in default browser",
	Run: func(cmd *cobra.Command, args []string) {
		spfr, _ := util.DetectSpfr(args)
		open(spfr)
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}

func open(spfr string) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(arg.Contest) == 0 {
		color.Red("Contest id not specified")
		os.Exit(1)
	}

	fmt.Println(color.BlueString("Opening problem's page:"), arg)
	util.BrowserOpen(arg.ProblemsPage())
}
