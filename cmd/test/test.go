package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

	"github.com/fatih/color"
)

// Test tests
func Test(submissionFilePath, checkerScript string,
	timeLimit time.Duration, memoryLimit uint64,
	inputStreamFile, outputStreamFile string,
	testingMode string, cnf *conf.Conf) {

	templateAlias := SelectSubmissionFile(&submissionFilePath, cnf)

	// Create temp directory for pseudo-sandbox testing.
	sandboxDir, err := ioutil.TempDir(os.TempDir(), "cpt-test-sandbox-")
	if err != nil {
		fmt.Println(color.RedString("error while creating sandbox:"), err)
		os.Exit(1)
	}

	// Copy submission file to sandbox directory.
	submissionFilePath, _ = filepath.Abs(submissionFilePath)
	if err := os.Symlink(submissionFilePath, filepath.Join(sandboxDir, filepath.Base(submissionFilePath))); err != nil {
		fmt.Println(color.RedString("error while moving submission file to sandbox:"), err)
		os.Exit(1)
	}

	// Manage placeholders for all test script template data.
	testScriptTmpltData := map[string]interface{}{
		"file":              filepath.Join(sandboxDir, filepath.Base(submissionFilePath)),
		"fileBasename":      filepath.Base(submissionFilePath),
		"fileBasenameNoExt": filepath.Ext(filepath.Base(submissionFilePath)),

		"sandboxDir": sandboxDir,
	}

	if preScript := cnf.GetString("template." + templateAlias + ".preScript"); preScript != "" {
		preScript, err := utils.CleanTemplate(preScript, testScriptTmpltData)
		if err != nil {
			fmt.Println(color.RedString("error while parsing prescript:"), err)
			os.Exit(1)
		}

		fmt.Println(color.BlueString("prescript:"), preScript)
		fmt.Println()

		if _, _, err := Execute(sandboxDir, preScript, nil, os.Stdout, os.Stderr, time.Minute, 256*1024*1024); err != nil {
			fmt.Println(color.RedString("error running prescript:"), err)
			os.Exit(1)
		}
	}

	if runScript := cnf.GetString("template." + templateAlias + ".runScript"); runScript != "" {
		runScript, err := utils.CleanTemplate(runScript, testScriptTmpltData)
		if err != nil {
			fmt.Println(color.RedString("error while parsing runscript:"), err)
			os.Exit(1)
		}

		// Available testing modes are:
		// d => default: Runs on each specified test case
		//      from the provided configuration.
		//
		// c => custom: Starts an interactive prompt
		//      with stdin/stdout to the terminal.
		//
		// i => interactive: Runs two commands, connecting
		//      stdin/stdout of first command to stdout/stdin
		//      of second command. (TODO)

		switch testingMode {
		case "d":
			// Default testing mode.

			// Extract list of testcase input/output files.
			testInputFiles := cnf.GetStrings("problem.test.input")
			testOutputFiles := cnf.GetStrings("problem.test.output")
			if len(testInputFiles) != len(testOutputFiles) {
				fmt.Println(color.RedString("error parsing testcase files:"),
					fmt.Sprintf("count of input files [%d] not equals count of output files [%d]",
						len(testInputFiles), len(testOutputFiles)))
				os.Exit(1)
			}

			for testIndex := 0; testIndex < len(testInputFiles); testIndex++ {
				// Create tests directory (to store current test files).
				testDir, _ := ioutil.TempDir(sandboxDir, "tests-")

				// Notice the slight difference in variable naming.
				// Here, 'output' and 'expected' notations have been swapped,
				// to stay in line with defaultTestingMode, testlib
				// and for backward compability.

				testInputFile := filepath.Join(testDir, "input.txt")
				testOutputFile := filepath.Join(testDir, "output.txt")
				testExpectedFile := filepath.Join(testDir, "expected.txt")

				var verd testExecDetails

				// Copy input and expected files. Create output file.
				testInputFiles[testIndex], _ = filepath.Abs(testInputFiles[testIndex])
				if err := os.Symlink(testInputFiles[testIndex], testInputFile); err != nil {
					verd = testExecDetails{
						verdict: color.HiRedString("FAIL"),
						failLog: err.Error(),
					}
					verd.prettyPrint(testIndex)
					continue
				}
				testOutputFiles[testIndex], _ = filepath.Abs(testOutputFiles[testIndex])
				if err := os.Symlink(testOutputFiles[testIndex], testExpectedFile); err != nil {
					verd = testExecDetails{
						verdict: color.HiRedString("FAIL"),
						failLog: err.Error(),
					}
					verd.prettyPrint(testIndex)
					continue
				}
				if _, err := os.Create(testOutputFile); err != nil {
					verd = testExecDetails{
						verdict: color.HiRedString("FAIL"),
						failLog: err.Error(),
					}
					verd.prettyPrint(testIndex)
					continue
				}

				verd = defaultTestingMode(sandboxDir, runScript, checkerScript,
					testInputFile, testOutputFile, testExpectedFile,
					inputStreamFile, outputStreamFile,
					timeLimit, memoryLimit)

				verd.prettyPrint(testIndex)
			}

		case "c":
			// Custom testing mode.

			customTestingMode(sandboxDir, runScript,
				timeLimit, memoryLimit)
		}
	}
}
