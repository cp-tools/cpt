package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/cp-tools/cpt-lib/codeforces"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
)

func contestsMode(arg codeforces.Args, count uint) {
	// determine number of pages to parse.
	pageCount := (count-1)/50 + 1

	chanContests, err := arg.GetContests(pageCount)
	if err != nil {
		fmt.Println(color.RedString("error while fetching contest details:"), err)
		os.Exit(1)
	}
	// Temporary color to prettify headers of table.
	cc := color.New(color.FgBlue, color.Underline).SprintFunc()

	// Create table and set rows headers.
	t := uitable.New()
	t.Separator = " | "
	t.MaxColWidth = 30
	t.Wrap = true
	t.AddRow(cc("ID"), cc("NAME"), cc("WRITERS"), cc("TIMINGS"), cc("REGISTRATION"))

	// Set live updater writer.
	writer := uilive.New()
	writer.Start()

	for contests := range chanContests {
		for _, contest := range contests {
			// We have to only print count rows of data.
			if count == 0 {
				break
			}
			// Pretty format timings data.
			var timings string
			timings += fmt.Sprintf("Begins: %v\n", contest.StartTime.Local().Format("Jan/02/2006 15:04"))
			timings += fmt.Sprintf("Length: %v\n", contest.Duration.String())

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

			t.AddRow(
				contest.Arg.Contest,                 // Contest ID
				contest.Name,                        // Contest name
				strings.Join(contest.Writers, ", "), // Contest writers
				timings,                             // Contest timings
				registrationStatus,                  // Registration status
			)
			// Added one more row to the table. Decrease count.
			count--
		}
		fmt.Fprintln(writer, t.String())
	}

	writer.Stop()
}
