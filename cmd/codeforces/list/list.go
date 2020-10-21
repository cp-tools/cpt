package list

import "github.com/cp-tools/cpt-lib/codeforces"

// List lists
func List(arg codeforces.Args, mode, username string, count uint) {
	switch mode {
	case "c": // contests
		contestsMode(arg, count)
	case "d": // dashboard
	case "s": // submissions
	}
}
