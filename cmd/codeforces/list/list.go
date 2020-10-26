package list

import (
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"

	"github.com/olekukonko/tablewriter"
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
func CompressVerdicts(verdict string) (string, tablewriter.Colors) {
	ccMap := [][]interface{}{
		{"Accepted", "Accepted", tablewriter.FgHiGreenColor},
		{"Partial result", "PR", tablewriter.FgHiGreenColor},
		{"Compilation error", "CE", tablewriter.FgHiYellowColor},
		{"Wrong answer", "WA", tablewriter.FgHiRedColor},
		{"Runtime error", "RTE", tablewriter.FgHiRedColor},
		{"Time limit exceeded", "TLE", tablewriter.FgHiYellowColor},
	}

	for _, val := range ccMap {
		if key := val[0].(string); strings.Contains(verdict, key) {
			// Clean text; update color; return
			verdict = strings.ReplaceAll(verdict, key, val[1].(string))
			return verdict, tablewriter.Colors{val[2].(int)}
		}
	}
	return verdict, tablewriter.Colors{}
}
