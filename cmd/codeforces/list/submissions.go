package list

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Submissions displays tabular submission data.
func Submissions(arg codeforces.Args, username string, count uint) {
	pageCount := (count-1)/50 + 1
	chanSubmissions, err := arg.GetSubmissions(username, pageCount)
	if err != nil {
		fmt.Println(color.RedString("error while fetching submission details:"), err)
		os.Exit(1)
	}

	// Set live updater writer.
	writer := uilive.New()
	writer.Start()

	// Create table to use.
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false

	headerColor := text.Colors{text.FgBlue, text.Bold}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignRight, WidthMax: 11},
		{Number: 2, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignLeft, WidthMax: 30},
		{Number: 3, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignCenter, WidthMax: 15},
		{Number: 4, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignCenter, WidthMax: 15},
		{Number: 5, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignRight, WidthMax: 12},
		{Number: 6, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignRight, WidthMax: 12},
	})

	t.AppendHeader(table.Row{"ID", "PROBLEM", "LANG", "VERDICT", "TIME", "MEMORY"})

	for submissions := range chanSubmissions {
		for i, submission := range submissions {
			// We have to only print count rows of data.
			if (pageCount == 1 && uint(i) >= count) ||
				(pageCount > 1 && count == 0) {
				break
			}

			verdict := ColorVerdict(submission)

			t.AppendRow(table.Row{
				submission.ID,       // Submission ID
				submission.Problem,  // Problem name
				submission.Language, // Submission language
				verdict,             // Submission verdict
				submission.Time,     // Time consumed
				submission.Memory,   // Memory consumed
			})

			if pageCount > 1 {
				count--
			}
		}
		if pageCount == 1 {
			// Continuous rendering when pageCount
			// is 1. Else, render all rows at once,
			// after all required rows are parsed.
			fmt.Fprintln(writer, t.Render())
			// Clear the table and add the header (again).
			t.ResetRows()
		}
	}
	if pageCount > 1 {
		fmt.Fprintln(writer, t.Render())
	}

	writer.Stop()
}

// ColorVerdict returns color coded verdict of submission.
func ColorVerdict(sub codeforces.Submission) string {
	ccMap := map[int]color.Attribute{
		codeforces.VerdictAC:          color.FgHiGreen,
		codeforces.VerdictPretestPass: color.FgHiGreen,

		codeforces.VerdictWA:   color.FgHiRed,
		codeforces.VerdictRTE:  color.FgHiRed,
		codeforces.VerdictDOJ:  color.FgHiRed,
		codeforces.VerdictHack: color.FgHiRed,

		codeforces.VerdictCE:  color.FgHiYellow,
		codeforces.VerdictTLE: color.FgHiYellow,
		codeforces.VerdictMLE: color.FgHiYellow,
		codeforces.VerdictILE: color.FgHiYellow,

		codeforces.VerdictSkip: color.FgHiCyan,
	}

	if col, ok := ccMap[sub.VerdictStatus]; ok {
		return color.New(col).Sprint(sub.Verdict)
	}

	return sub.Verdict
}
