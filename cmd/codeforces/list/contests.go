package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Contests displays tabular contest data.
func Contests(arg codeforces.Args, count uint) {
	pageCount := (count-1)/100 + 1
	chanContests, err := arg.GetContests(pageCount)
	if err != nil {
		fmt.Println(color.RedString("error while fetching contest details:"), err)
		os.Exit(1)
	}

	// Create table to use.
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false

	headerColor := text.Colors{text.FgBlue, text.Bold}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignRight, WidthMax: 8},
		{Number: 2, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignLeft, WidthMax: 30},
		{Number: 3, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignCenter, WidthMax: 25},
		{Number: 4, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignCenter, WidthMax: 30},
		{Number: 5, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, Align: text.AlignCenter, WidthMax: 14},
	})
	t.Style().Options.SeparateRows = true

	t.AppendHeader(table.Row{"ID", "NAME", "WRITERS", "TIMINGS", "REGISTRATION"})

	for contests := range chanContests {
		for _, contest := range contests {
			// We have to only print count rows of data.
			if count == 0 {
				break
			}
			// Pretty format timings data.
			var timings string
			timings += fmt.Sprintf("Begins: %v\n", contest.StartTime.Local().Format("Jan/02/2006 15:04"))
			timings += fmt.Sprintf("Length: %v", contest.Duration.String())

			// @todo Hyperlink registration status text
			// @body Make hyperlink to registration page if registration is open.
			// @body Also, add support for virtual registration.

			var registrationStatus string
			switch contest.RegStatus {
			case codeforces.RegistrationOpen:
				registrationStatus = color.HiGreenString("OPEN")
			case codeforces.RegistrationClosed:
				registrationStatus = color.HiRedString("CLOSED")
			case codeforces.RegistrationDone:
				registrationStatus = color.HiGreenString("DONE")
			case codeforces.RegistrationNotExists:
				registrationStatus = color.HiYellowString("NA")
			}

			t.AppendRow(table.Row{
				contest.Arg.Contest,                // Contest ID
				contest.Name,                       // Contest name
				strings.Join(contest.Writers, " "), // Contest writers
				timings,                            // Contest timings
				registrationStatus,                 // Registration status
			})

			count--
		}
	}
	fmt.Println(t.Render())
}
