package test

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/cp-tools/cpt/packages/conf"

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
	// Load all scripts into template (check if correctly formed).
	preScript := cnf.GetString("template." + alias + ".preScript")
	runScript := cnf.GetString("template." + alias + ".runScript")
	postScript := cnf.GetString("template." + alias + ".postScript")

	// Run preScript.
	if preScript != "" {
		var script strings.Builder
		tmplt := template.Must(template.New("").Parse(preScript))
		tmplt.Execute(&script, tmpltData)
		fmt.Println(color.BlueString("prescript:"), script.String())

		if _, err := runShellScript(script.String(), time.Minute,
			os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println()
	}

	// Run script for tests.
	if runScript != "" {
		var script strings.Builder
		tmplt := template.Must(template.New("").Parse(runScript))
		tmplt.Execute(&script, tmpltData)

		switch mode {
		case "j": // Default judge mode.
			inputFiles, expectedFiles := extractTestsFiles(cnf)
			for i := 0; i < len(inputFiles); i++ {
				judgeMode(script.String(), timelimit, inputFiles[i], expectedFiles[i], checker, i)
			}

		case "i": // Interactive mode.
			interactiveMode(script.String())
		}
	}

	// Run postScript.
	if postScript != "" {
		var script strings.Builder
		tmplt := template.Must(template.New("").Parse(postScript))
		tmplt.Execute(&script, tmpltData)
		fmt.Println(color.BlueString("postscript:"), script.String())

		if _, err := runShellScript(script.String(), time.Minute,
			os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
