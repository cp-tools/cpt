package list

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Dashboard displays tabulated dashboard data.
func Dashboard(arg codeforces.Args) {
	dashboard, err := arg.GetDashboard()
	if err != nil {
		fmt.Println(color.RedString("error while fetching dashboard details:"), err)
		os.Exit(1)
	}

	fmt.Println(color.BlueString("Contest name:"), dashboard.Name)
	if dashboard.Countdown > 0 {
		fmt.Println(color.BlueString("Contest ends in:"), dashboard.Countdown)
	}
	fmt.Println()

	// Create table to use.
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false

	headerColor := text.Colors{text.FgBlue, text.Bold}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignRight, WidthMax: 8},
		{Number: 2, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignLeft, WidthMax: 30},
		{Number: 3, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignCenter, WidthMax: 10},
		{Number: 4, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignRight, WidthMax: 10},
	})

	t.AppendHeader(table.Row{"#", "NAME", "STATUS", "SOLVED"})

	for _, problem := range dashboard.Problem {
		status := ""
		switch problem.SolveStatus {
		case codeforces.SolveAccepted:
			status = color.HiGreenString("AC")
		case codeforces.SolveRejected:
			status = color.HiRedString("WA")
		case codeforces.SolveNotAttempted:
			status = "NA"
		}

		t.AppendRow(table.Row{
			problem.Arg.Problem, // Problem ID
			problem.Name,        // Problem name
			status,              // Solved status
			problem.SolveCount,  // Solve count
		})

	}
	fmt.Println(t.Render())
}
