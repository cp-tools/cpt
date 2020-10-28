package list

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
)

// DashboardMode displays tabulated dashboard data.
func DashboardMode(arg codeforces.Args) {
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
