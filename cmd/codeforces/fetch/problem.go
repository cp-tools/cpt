package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

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
	// create tests/ folder, to save tests to.
	if err := os.MkdirAll("tests", os.ModePerm); err != nil {
		fmt.Println(color.RedString("error occurred while creating tests folder:"), err)
		os.Exit(1)
	}

	lastTestIndex := 0
	for _, sampleTest := range problem.SampleTests {
		inFileName := func() string { return filepath.Join("tests", fmt.Sprintf("%d.in", lastTestIndex)) }
		outFileName := func() string { return filepath.Join("tests", fmt.Sprintf("%d.out", lastTestIndex)) }

		for utils.FileExists(inFileName()) || utils.FileExists(outFileName()) {
			lastTestIndex++
		}

		inFile, err := os.Create(inFileName())
		if err != nil {
			fmt.Println(color.RedString("error while creating file:"), err)
			return
		}
		outFile, err := os.Create(outFileName())
		if err != nil {
			fmt.Println(color.RedString("error while creating file:"), err)
			return
		}

		// Write data to respective files.
		inFile.WriteString(sampleTest.Input)
		outFile.WriteString(sampleTest.Output)
		// Append names to slices.
		inFiles = append(inFiles, inFileName())
		outFiles = append(outFiles, outFileName())

		inFile.Close()
		outFile.Close()
	}
	// Verbose message regarding how many tests parsed.
	cnt := len(problem.SampleTests)
	gs := color.New(color.FgGreen).SprintFunc()
	fmt.Println(gs("fetched"), cnt, gs("sample tests in"), problem.Name)
	return
}
