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

func contestsMode(arg codeforces.Args, count uint) {
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

func dashboardMode(arg codeforces.Args) {
	dashboard, err := arg.GetDashboard()
	if err != nil {
		fmt.Println(color.RedString("error while fetching dashboard details:"), err)
		os.Exit(1)
	}

	fmt.Println(color.BlueString("Contest name:"), dashboard.Name)
	if dashboard.Countdown > 0 {
		fmt.Println(color.BlueString("Contest ends in:"), dashboard.Countdown)
	}

	t := uitable.New()
	t.Separator = " | "
	t.MaxColWidth = 40
	t.Wrap = true

	hdr := util.ColorHeaderFormat("#", "NAME", "STATUS", "SOLVED")
	t.AddRow(hdr[0], hdr[1], hdr[2], hdr[3])

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

		t.AddRow(
			problem.Arg.Problem, // Problem ID
			problem.Name,        // Problem name
			status,              // Solved status
			problem.SolveCount,  // Solve count
		)

	}
	fmt.Println(t.String())
}

func submissionsMode(arg codeforces.Args, handle string, count uint) {
	// Anything more than 1 page (50 rows) makes no sense.
	chanSubmissions, err := arg.GetSubmissions(handle, 1)
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
