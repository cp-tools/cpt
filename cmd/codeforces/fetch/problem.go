package fetch

import (
	"fmt"
	"os"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/packages/conf"
	"github.com/fatih/color"
)

func createConf(problem codeforces.Problem, testInFiles, testOutFiles []string) *conf.Conf {
	cnf := conf.New("problem")
	cnf.LoadFile("meta.yaml")
	/*
		problem:
			name: XXXX
			timeLimit: XXXX
			memoryLimit: XXXX
			arg:
				Contest: XXXX
				Problem: XXXX
				Class: XXXX
				Group: XXXX
			inputStream: XXXX
			outputStream: XXXX
			tests:
				input: [XXXX]
				output: [XXXX]
	*/
	cnf.Set("problem.name", problem.Name)
	cnf.Set("problem.timeLimit", problem.TimeLimit)
	cnf.Set("problem.memoryLimit", problem.MemoryLimit)
	cnf.Set("problem.arg", problem.Arg)
	cnf.Set("problem.inputStream", problem.InpStream)
	cnf.Set("problem.outputStream", problem.OutStream)
	cnf.Set("problem.test.input", testInFiles)
	cnf.Set("problem.test.output", testOutFiles)
	cnf.WriteFile()
	return cnf
}

func createTests(problem codeforces.Problem) (inFiles, outFiles []string) {
	for i := range problem.SampleTests {
		inFile, err := os.Create(fmt.Sprintf("%d.in", i))
		if err != nil {
			fmt.Println(color.RedString("error while creating file:"), err)
			return
		}
		outFile, err := os.Create(fmt.Sprintf("%d.out", i))
		if err != nil {
			fmt.Println(color.RedString("error while creating file:"), err)
			return
		}

		// Write data to respective files.
		inFile.WriteString(problem.SampleTests[i].Input)
		outFile.WriteString(problem.SampleTests[i].Output)
		// Append names to slices.
		inFiles = append(inFiles, inFile.Name())
		outFiles = append(outFiles, outFile.Name())

		inFile.Close()
		outFile.Close()
	}
	// Verbose message regarding how many tests parsed.
	cnt := len(problem.SampleTests)
	gs := color.New(color.FgGreen).SprintFunc()
	fmt.Println(gs("fetched"), cnt, gs("sample tests in"), problem.Name)
	return
}
