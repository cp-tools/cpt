package submit

import (
	"fmt"
	"os"
	"strings"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/cmd/test"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/olekukonko/tablewriter"
)

// Submit submits
// filePath must point to a valid file.
func Submit(arg codeforces.Args, filePath string, cnf *conf.Conf) {
	// Determine code file and template alias to use.
	fileName, alias := test.SelectCodeFile(filePath, cnf)
	// Submit solution.
	langName := cnf.GetString("template." + alias + ".language")
	submission, err := arg.SubmitSolution(langName, fileName)
	if err != nil {
		fmt.Println(color.RedString("error submitting solution:"), err)
		os.Exit(1)
	}
	fmt.Println(color.GreenString("submitted successfully"))

	// Table to use to display verdict.
	tString := &strings.Builder{}
	t := tablewriter.NewWriter(tString)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.SetTablePadding("\t")
	t.SetBorder(false)

	// Run live verdict till judging completed.
	writer := uilive.New()
	writer.Start()
	for sub := range submission {
		t.ClearRows()
		tString.Reset()

		t.Append("VERDICT:", sub.Verdict)
		if sub.IsJudging == false {
			// Judging done; add resource data.
			t.Append("MEMORY:", sub.Memory)
			t.Append("TIME:", sub.Time)
		}

		t.Render()
		fmt.Fprintln(writer, tString.String())
	}
	writer.Stop()
}
