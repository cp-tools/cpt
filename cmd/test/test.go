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
func Test(checker, filePath string, interactive bool, timelimit time.Duration, cnf *conf.Conf) {
	// Determine code file and template alias to use.
	//fileName, alias := SelectCodeFile(filePath, cnf)
	_, alias := SelectCodeFile(filePath, cnf)
	// Configure all template placeholder fields here.
	tmpltData := map[string]interface{}{
		"file": filePath,
	}
	// Load all scripts into template (check if correctly formed).
	preScript := cnf.GetString("template." + alias + ".preScript")
	//runScript := cnf.GetString("template." + alias + ".runScript")
	postScript := cnf.GetString("template." + alias + ".postScript")

	// Run preScript.
	if preScript != "" {
		var script strings.Builder
		tmplt := template.Must(template.New("").Parse(preScript))
		tmplt.Execute(&script, tmpltData)
		fmt.Println(color.BlueString("prescript:"), script)

		if _, err := runShellScript(script.String(), time.Minute,
			os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	// Run postScript.
	if postScript != "" {
		var script strings.Builder
		tmplt := template.Must(template.New("").Parse(postScript))
		tmplt.Execute(&script, tmpltData)
		fmt.Println(color.BlueString("postscript:"), script)

		if _, err := runShellScript(script.String(), time.Minute,
			os.Stdin, os.Stdout, os.Stderr); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
