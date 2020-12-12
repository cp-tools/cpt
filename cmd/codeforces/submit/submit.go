package submit

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/cmd/codeforces/list"
	"github.com/cp-tools/cpt/cmd/test"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/jedib0t/go-pretty/v6/table"
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
	fmt.Println(color.GreenString("Submitted successfully!"))

	// Table to use to display verdict.
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false
	t.Style().Box.PaddingRight = "\t"

	// Run live verdict till judging completed.
	writer := uilive.New()
	writer.Start()

	for sub := range submission {
		t.ResetRows()

		t.AppendRow(table.Row{color.BlueString("Verdict:"), list.ColorVerdict(sub)})
		if sub.IsJudging == false {
			// Judging done; add resource data.
			t.AppendRow(table.Row{color.BlueString("Memory:"), sub.Memory})
			t.AppendRow(table.Row{color.BlueString("Time:"), sub.Time})
		}

		fmt.Fprintln(writer, t.Render())
	}

	writer.Stop()
}
