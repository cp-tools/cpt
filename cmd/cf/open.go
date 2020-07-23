package cf

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"
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
		fmt.Println("Contest id not specified")
		os.Exit(1)
	}

	fmt.Println("Opening problem's page:", arg.Contest, arg.Problem)
	util.BrowserOpen(arg.ProblemsPage())
}
