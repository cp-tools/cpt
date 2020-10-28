package submit

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/cmd/codeforces/list"
	"github.com/cp-tools/cpt/cmd/test"
	"github.com/cp-tools/cpt/packages/conf"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
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
	t := uitable.New()
	t.Separator = "\t"

	// Run live verdict till judging completed.
	writer := uilive.New()
	writer.Start()

	for sub := range submission {
		t.Rows = nil

		verdict := list.CompressVerdicts(sub.Verdict)

		t.AddRow(util.ColorSetBlueBold("Verdict:"), verdict)
		if sub.IsJudging == false {
			// Judging done; add resource data.
			t.AddRow(util.ColorSetBlueBold("Memory:"), sub.Memory)
			t.AddRow(util.ColorSetBlueBold("Time:"), sub.Time)
		}

		fmt.Fprintln(writer, t.String())
	}

	writer.Stop()
}
