package test

import (
	"fmt"
	"os"
	"time"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
)

// Test tests
func Test(checker, filePath, mode string, timelimit time.Duration, cnf *conf.Conf) {
	// Determine code file and template alias to use.
	fileName, alias := SelectCodeFile(filePath, cnf)
	// Configure all template placeholder fields here.
	tmpltData := map[string]interface{}{
		"file": fileName,
	}

	// Run preScript.
	if preScript := cnf.GetString("template." + alias + ".preScript"); preScript != "" {
		script, _ := util.CleanTemplate(preScript, tmpltData)
		fmt.Println(color.BlueString("prescript:"), script, "\n")

		if _, err := runShellScript(script, time.Minute, os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	// Run script for tests.
	if runScript := cnf.GetString("template." + alias + ".runScript"); runScript != "" {
		script, _ := util.CleanTemplate(runScript, tmpltData)

		switch mode {
		case "j": // Default judge mode.
			checkerTmplt := cnf.GetString("checker.checkers." + checker + ".script")

			inputFiles, expectedFiles := extractTestsFiles(cnf)
			for i := 0; i < len(inputFiles); i++ {
				judgeMode(script, checkerTmplt, timelimit, inputFiles[i], expectedFiles[i], i)
				if i != len(inputFiles)-1 {
					// Print newline after every (but last) test case.
					fmt.Println()
				}
			}

		case "i": // Interactive mode.
			interactiveMode(script)
		}
	}

	// Run postScript.
	if postScript := cnf.GetString("template." + alias + ".postScript"); postScript != "" {
		script, _ := util.CleanTemplate(postScript, tmpltData)
		fmt.Println("\n", color.BlueString("postscript:"), script)

		if _, err := runShellScript(script, time.Minute, os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
