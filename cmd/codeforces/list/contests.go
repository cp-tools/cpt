package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
)

// Contests displays tabular contest data.
func Contests(arg codeforces.Args, count uint) {
	// Anything more than 1 page (100 rows) makes no sense.
	chanContests, err := arg.GetContests(1)
	if err != nil {
		fmt.Println(color.RedString("error while fetching contest details:"), err)
		os.Exit(1)
	}

	// Create table to use.
	t := uitable.New()
	t.Separator = " | "
	t.MaxColWidth = 30
	t.Wrap = true

	hdr := util.ColorHeaderFormat("ID", "NAME", "WRITERS", "TIMINGS", "REGISTRATION")
	t.AddRow(hdr[0], hdr[1], hdr[2], hdr[3], hdr[4])

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

			count--
		}
	}
	fmt.Println(t.String())
}
