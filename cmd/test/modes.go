package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func judgeMode(script string, timelimit time.Duration, inputFile, expectedFile, checkerPath string, index int) {
	// A hacky function to color certain parts of a template.
	c := color.New(color.FgBlue, color.Bold).SprintFunc()
	// Verdict template configurations.
	tmpltData := map[string]interface{}{}
	tmplt := template.Must(template.New("").Parse(
		c("Test:") + " #{{.index}}    " + c("Verdict:") + " {{.verdict}}    " + c("Time:") + " {{.elapsed}}\n" +
			"{{- if .failLog}}\n" + c("Fail:") + "\n{{.failLog}}{{end}}\n" +
			"{{- if .stderr}}\n" + c("Stderr:") + "\n{{.stderr}}{{end}}\n" +
			"{{- if .checkerLog}}\n" + c("Checker Log:") + " {{.checkerLog}}{{end}}\n" +
			"{{- if .testDetails}}\n" + c("Input:") + "\n{{.input}}\n{{.testDetails}}{{end}}\n",
	))

	defer func() {
		tmpltData["index"] = index
		// handle panic; recover.
		if r := recover(); r != nil {
			tmpltData["verdict"] = color.RedString("FAIL")
			tmpltData["failLog"] = r.(error)
		}
		// Print verdict data to stdout.
		tmplt.Execute(os.Stdout, tmpltData)
	}()

	// Read input from file.
	input, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	// Run code against test case.
	var output, stderr bytes.Buffer
	elapsed, err := runShellScript(script, timelimit, input, &output, &stderr)
	// Common to all template values to write.
	tmpltData["elapsed"] = elapsed.Truncate(time.Millisecond)
	tmpltData["stderr"] = stderr.String()

	// Determine verdicts.
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		// It's a time limit exceeded case.
		tmpltData["verdict"] = color.YellowString("TLE")

	case err != nil:
		// It's a runtime error case.
		tmpltData["verdict"] = color.RedString("RTE")

	default:
		// Output recieved, and program exited normally.
		// Check if output matches required answer(s).

		// write output to temporary file.
		outputFile, err := ioutil.TempFile(os.TempDir(), "cpt-test-output-")
		if err != nil {
			panic(err)
		}
		defer os.Remove(outputFile.Name())
		outputFile.Write(output.Bytes())

		// Run checker to validate output. Shell script to run is:
		// (<checker> <input-file> <expected-file> <recieved-file>)
		var checkerStderr bytes.Buffer
		checkerScript := fmt.Sprintf("%v %v %v %v", checkerPath, inputFile, outputFile.Name(), expectedFile)
		_, err = runShellScript(checkerScript, time.Minute, nil, nil, &checkerStderr)
		// Set template field data.
		tmpltData["checkerLog"] = checkerStderr.String()

		if _, ok := err.(*exec.ExitError); ok {
			// Checker ended with non-zero error code.
			// Verdict is thus wrong answer.
			tmpltData["verdict"] = color.RedString("WA")

			// Read input from file.
			input, err := os.Open(inputFile)
			if err != nil {
				panic(err)
			}
			defer input.Close()

			inputBuf := make([]byte, 80)
			if n, _ := input.Read(inputBuf); n == len(inputBuf) {
				inputBuf = append(inputBuf[:n-3], []byte("...")...)
			}

			// Read expected from file.
			expected, err := os.Open(expectedFile)
			if err != nil {
				panic(err)
			}
			defer expected.Close()

			expectedBuf := make([]byte, 80)
			if n, _ := expected.Read(expectedBuf); n == len(expectedBuf) {
				expectedBuf = append(expectedBuf[:n-3], []byte("...")...)
			}

			// Read output from created output file.
			output, err := os.Open(outputFile.Name())
			if err != nil {
				panic(err)
			}
			defer output.Close()

			outputBuf := make([]byte, 80)
			if n, _ := output.Read(outputBuf); n == len(outputBuf) {
				outputBuf = append(outputBuf[:n-3], []byte("...")...)
			}

			// Temporary color to prettify headers of table.
			tHeaderCol := tablewriter.Color(tablewriter.FgBlueColor, tablewriter.Bold)
			// Table to display output difference.
			tString := &strings.Builder{}
			t := tablewriter.NewWriter(tString)
			t.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
			t.SetHeaderColor(tHeaderCol, tHeaderCol)
			t.SetCenterSeparator("")
			t.SetColumnSeparator("")
			t.SetRowSeparator("")
			t.SetTablePadding("\t")
			t.SetBorder(false)
			t.SetColWidth(50)

			t.Append("OUTPUT", "EXPECTED")
			t.Append(string(outputBuf), string(expectedBuf))

			t.Render()
			tmpltData["testDetails"] = tString.String()
			tmpltData["input"] = string(inputBuf)

		} else if err != nil {
			// Unknown error; Panic.
			panic(err)
		} else {
			// Solution produced correct answer.
			tmpltData["verdict"] = color.GreenString("AC")
		}
	}
}

func interactiveMode(script string) {
	// It doesn't get any simpler, does it?
	fmt.Println(color.GreenString("---- * ---- launched ---- * ----"))
	runShellScript(script, time.Hour, os.Stdin, os.Stdout, os.Stderr)
	fmt.Println(color.GreenString("---- * ---- finished ---- * ----"))
	fmt.Println() // Newline for asthetics.
}
