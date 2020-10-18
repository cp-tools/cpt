package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
)

func judgeMode(script string, timelimit time.Duration, inputFile, expectedFile, checkerPath string, index int) {
	// A hacky function to color certain parts of a template.
	c := color.New(color.FgBlue, color.Bold).SprintFunc()
	// Verdict template configurations.
	tmpltData := map[string]interface{}{}
	tmplt := template.Must(template.New("").Parse(
		c("Test:") + " #{{.index}}\t" + c("Verdict:") + " {{.verdict}}\t" + c("Time:") + " {{.elapsed}}\n" +
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
		defer outputFile.Close()
		outputFile.WriteString(output.String())

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
			if _, err := input.Read(inputBuf); err == nil {
				inputBuf = append(inputBuf, []byte("...")...)
			}

			// Read expected from file.
			expected, err := os.Open(expectedFile)
			if err != nil {
				panic(err)
			}
			defer expected.Close()

			expectedBuf := make([]byte, 80)
			if _, err := expected.Read(expectedBuf); err == nil {
				expectedBuf = append(expectedBuf, []byte("...")...)
			}

			// Read output from created output file.
			outputBuf := make([]byte, 80)
			if _, err := outputFile.Read(outputBuf); err == nil {
				outputBuf = append(outputBuf, []byte("...")...)
			}

			testDetails := uitable.New()
			testDetails.Separator = "\t|"
			testDetails.AddRow(c("OUTPUT:"), c("EXPECTED:"))
			testDetails.AddRow(string(outputBuf), string(expectedBuf))

			tmpltData["testDetails"] = testDetails.String()
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
