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

func judgeMode(script, checkerTmplt string, timelimit time.Duration,
	inputFile, expectedFile string, index int) {
	// A hacky function to color certain parts of a template.
	c := color.New(color.FgBlue, color.Bold).SprintFunc()

	// Verdict template configurations.
	verdictTmpltData := map[string]interface{}{}
	tmplt := template.Must(template.New("verdict").Parse(
		c("Test:") + " #{{.index}}    " + c("Verdict:") + " {{.verdict}}    " + c("Time:") + " {{.elapsed}}\n" +
			"{{- if .failLog}}\n" + c("Fail:") + "{{.failLog}}{{end}}\n" +
			"{{- if .stderr}}\n" + c("Stderr:") + "\n{{.stderr}}{{end}}\n" +
			"{{- if .checkerLog}}\n" + c("Checker Log:") + " {{.checkerLog}}{{end}}\n" +
			"{{- if .testDetails}}\n" + c("Input:") + "\n{{.input}}\n{{.testDetails}}{{end}}\n",
	))

	// Checker template configurations.
	checkerTmpltData := map[string]interface{}{
		"inputFile":    inputFile,
		"expectedFile": expectedFile,
	}
	template.Must(tmplt.New("checker").Parse(checkerTmplt))

	// handle panic.
	defer func() {
		verdictTmpltData["index"] = index
		// handle panic; recover.
		if r := recover(); r != nil {
			verdictTmpltData["verdict"] = color.RedString("FAIL")
			verdictTmpltData["failLog"] = r.(error)
		}
		// Print verdict data to stdout.
		tmplt.ExecuteTemplate(os.Stdout, "verdict", verdictTmpltData)
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
	verdictTmpltData["elapsed"] = elapsed.Truncate(time.Millisecond)
	verdictTmpltData["stderr"] = stderr.String()

	// Determine verdicts.
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		// It's a time limit exceeded case.
		verdictTmpltData["verdict"] = color.YellowString("TLE")

	case err != nil:
		// It's a runtime error case.
		verdictTmpltData["verdict"] = color.RedString("RTE")

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

		// Run checker to validate output.
		var checkerScript strings.Builder
		checkerTmpltData["outputFile"] = outputFile.Name()
		tmplt.ExecuteTemplate(&checkerScript, "checker", checkerTmpltData)

		var checkerStderr bytes.Buffer
		_, err = runShellScript(checkerScript.String(), time.Minute, nil, nil, &checkerStderr)

		// Set template field data.
		verdictTmpltData["checkerLog"] = checkerStderr.String()
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			// Checker ended with error code 1.
			// Verdict is thus wrong answer.
			verdictTmpltData["verdict"] = color.RedString("WA")

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
			t.SetCenterSeparator("")
			t.SetColumnSeparator("")
			t.SetRowSeparator("")
			t.SetNoWhiteSpace(true)
			t.SetTablePadding("\t")
			t.SetHeaderLine(false)
			t.SetBorder(false)

			t.SetHeader("OUTPUT", "EXPECTED")
			t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			t.SetHeaderColor(tHeaderCol, tHeaderCol)

			t.Append(string(outputBuf), string(expectedBuf))

			t.Render()
			verdictTmpltData["testDetails"] = tString.String()
			verdictTmpltData["input"] = string(inputBuf)

		} else if ok {
			// Unknown error; Panic.
			panic(err)
		} else {
			// Solution produced correct answer.
			verdictTmpltData["verdict"] = color.GreenString("AC")
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
