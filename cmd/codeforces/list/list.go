package list

import (
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"

	"github.com/fatih/color"
)

// List lists
func List(arg codeforces.Args, mode, username string, count uint) {
	switch mode {
	case "c": // contests
		contestsMode(arg, count)
	case "d": // dashboard
		dashboardMode(arg)
	case "s": // submissions
		submissionsMode(arg, username, count)
	}
}

// CompressVerdicts compresses and returns color to use.
func CompressVerdicts(verdict string) string {
	ccMap := [][]interface{}{
		{"Accepted", "Accepted", color.FgHiGreen},
		{"Partial result", "PR", color.FgHiGreen},
		{"Compilation error", "CE", color.FgHiYellow},
		{"Wrong answer", "WA", color.FgHiRed},
		{"Runtime error", "RTE", color.FgHiRed},
		{"Time limit exceeded", "TLE", color.FgHiYellow},
	}

	for _, val := range ccMap {
		if key := val[0].(string); strings.Contains(verdict, key) {
			// Clean text; update color; return
			verdict = strings.ReplaceAll(verdict, key, val[1].(string))
			return color.New(val[2].(color.Attribute)).Sprint(verdict)
		}
	}
	return verdict
}
