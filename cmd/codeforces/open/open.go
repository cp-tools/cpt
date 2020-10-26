package open

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/cp-tools/cpt-lib/v2/codeforces"

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
	openURL(pageURL)
}

// Attribution: https://stackoverflow.com/a/39324149
// openURL opens the specified URL in the default browser of the user.
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
