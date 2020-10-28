package list

import (
	"strings"

	"github.com/fatih/color"
)

// Refer respective mode files for their implementation.

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
