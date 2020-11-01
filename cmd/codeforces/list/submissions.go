package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
)

// Submissions displays tabular submission data.
func Submissions(arg codeforces.Args, username string, count uint) {
	// Anything more than 1 page (50 rows) makes no sense.
	chanSubmissions, err := arg.GetSubmissions(username, 1)
	if err != nil {
		fmt.Println(color.RedString("error while fetching submission details:"), err)
		os.Exit(1)
	}

	// Set live updater writer.
	writer := uilive.New()
	writer.Start()

	// Create table to use.
	t := uitable.New()
	t.Separator = " | "
	t.MaxColWidth = 30
	t.Wrap = true

	hdr := util.ColorHeaderFormat("ID", "PROBLEM", "LANG", "VERDICT", "TIME", "MEMORY")
	t.AddRow(hdr[0], hdr[1], hdr[2], hdr[3], hdr[4], hdr[5])

	for submissions := range chanSubmissions {
		for _, submission := range submissions {
			// We have to only print count rows of data.
			if count == 0 {
				continue
			}

			verdict := CompressVerdicts(submission.Verdict)

			t.AddRow(
				submission.ID,       // Submission ID
				submission.Problem,  // Problem name
				submission.Language, // Submission language
				verdict,             // Submission verdict
				submission.Time,     // Time consumed
				submission.Memory,   // Memory consumed
			)

			count--
		}
		fmt.Fprintln(writer, t.String())
	}
	writer.Stop()
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