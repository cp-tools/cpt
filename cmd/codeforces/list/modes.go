package list

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/olekukonko/tablewriter"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

func contestsMode(arg codeforces.Args, count uint) {
	// Anything more than 1 page (100 rows) makes no sense.
	chanContests, err := arg.GetContests(1)
	if err != nil {
		fmt.Println(color.RedString("error while fetching contest details:"), err)
		os.Exit(1)
	}

	// Set live updater writer.
	writer := uilive.New()

	// Create table to use.
	t := tablewriter.NewWriter(writer)
	t.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	t.SetCenterSeparator("|")
	t.SetBorder(false)
	t.SetColWidth(30)

	// Temporary color to prettify headers of table.
	col := tablewriter.Color(tablewriter.FgBlueColor, tablewriter.Bold)
	t.SetHeader("ID", "NAME", "WRITERS", "TIMINGS", "REGISTRATION")
	t.SetHeaderColor(col, col, col, col, col)

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

			t.Append(
				contest.Arg.Contest,                 // Contest ID
				contest.Name,                        // Contest name
				strings.Join(contest.Writers, ", "), // Contest writers
				timings,                             // Contest timings
				registrationStatus,                  // Registration status
			)
			// Added one more row to the table. Decrease count.
			count--
		}
	}
	t.Render()
	writer.Stop()
}

func dashboardMode(arg codeforces.Args) {
	dashboard, err := arg.GetDashboard()
	if err != nil {
		fmt.Println(color.RedString("error while fetching dashboard details:"), err)
		os.Exit(1)
	}

	// A hacky function to color certain parts of a template.
	c := color.New(color.FgBlue, color.Bold).SprintFunc()

	fmt.Println(c("Contest name:"), dashboard.Name)
	if dashboard.Countdown > 0 {
		fmt.Println(c("Contest ends in:"), dashboard.Countdown)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	t.SetCenterSeparator("|")
	t.SetColWidth(30)

	col := tablewriter.Color(tablewriter.FgBlueColor, tablewriter.Bold)
	t.SetHeader("#", "Name", "Status", "Solved")
	t.SetHeaderColor(col, col, col, col)
	t.SetColumnAlignment(tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_CENTER, tablewriter.ALIGN_RIGHT)

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

		t.Append(
			problem.Arg.Problem,
			problem.Name,
			status,
			strconv.Itoa(problem.SolveCount),
		)
	}
	fmt.Println()
	t.Render()
}
