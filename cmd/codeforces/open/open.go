package open

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/utils"

	"github.com/fatih/color"
)

// Open launches the required webpage in users default browser.
func Open(arg codeforces.Args, mode string) {
	var pageURL string
	var err error

	switch mode {
	case "p": // problem page
		pageURL, err = arg.ProblemsPage()
	case "d": // dashboard page
		pageURL, err = arg.DashboardPage()
	case "s": // submission page
		pageURL, err = arg.SubmissionsPage("")
	}

	if err != nil {
		fmt.Println(color.RedString("error while determining page url:"), err)
		os.Exit(1)
	}
	// Open the webpage.
	utils.OpenURL(pageURL)
}
