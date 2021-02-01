package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cp-tools/cpt/utils"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type (
	testExecDetails struct {
		verdict    string
		failLog    string
		checkerLog string
		stderrLog  string

		timeConsumed   time.Duration
		memoryConsumed uint64

		testDetails struct {
			truncInput    string
			truncOutput   string
			truncExpected string
		}
	}
)

func (verd *testExecDetails) prettyPrint(testIndex int) {
	c := color.New(color.FgBlue, color.Bold).SprintFunc()

	// Test: 4    Verdict: WA    Time: 32ms
	fmt.Println(c("Test:"), testIndex, "\t"+c("Verdict:"), verd.verdict)
	// Time: 32ms    Memory: 2000kb
	verd.timeConsumed = verd.timeConsumed.Truncate(time.Millisecond)
	memoryKB := fmt.Sprintf("%dkb", verd.memoryConsumed/1024)
	fmt.Println(c("Time:"), verd.timeConsumed, "\t"+c("Memory:"), memoryKB)

	if verd.failLog != "" {
		// Fail: Could not execute checker
		fmt.Println(c("Fail log:"), strings.TrimSpace(verd.failLog))
	}

	if verd.stderrLog != "" {
		// Stderr:
		// 1 2 3
		// a b c
		fmt.Println(c("Stderr:"), strings.TrimSpace(verd.stderrLog))
	}

	if verd.checkerLog != "" {
		// Checker Log: Wrong answer, expected 3, found 4.
		fmt.Println(c("Checker log:"), strings.TrimSpace(verd.checkerLog))
	}

	if strings.Contains(verd.verdict, "WA") {
		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.Style().Options.DrawBorder = false

		headerColor := text.Colors{text.FgBlue, text.Bold}
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, WidthMax: 40},
			{Number: 2, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, WidthMax: 40},
			{Number: 3, AlignHeader: text.AlignCenter, ColorsHeader: headerColor, WidthMax: 40},
		})

		t.AppendHeader(table.Row{"INPUT", "EXPECTED", "OUTPUT"})
		t.AppendRow(table.Row{verd.testDetails.truncInput,
			verd.testDetails.truncExpected, verd.testDetails.truncOutput})

		fmt.Println()
		fmt.Println(t.Render())
	}
	fmt.Println()
}

func defaultTestingMode(sandboxDir, runScript, checkerScript string,
	testInputFile, testOutputFile, testExpectedFile string,
	testInputStreamFile, testOutputStreamFile string,
	timeLimit time.Duration, memoryLimit uint64) (verd testExecDetails) {

	// Unexpected errors - handle panics.
	defer func() {
		if r := recover(); r != nil {
			verd.verdict = color.HiRedString("FAIL")
			verd.failLog = r.(error).Error()
		}
	}()

	var testStdin io.Reader
	var testStdout io.Writer
	var testStderr strings.Builder

	// If the submission code reads from stdin.
	if testInputStreamFile == "" {
		fl, err := os.Open(testInputFile)
		if err != nil {
			panic(err)
		}
		defer fl.Close()

		testStdin = fl
	} else {
		err := os.Symlink(testInputFile, filepath.Join(sandboxDir, testInputStreamFile))
		if err != nil {
			panic(err)
		}
		testStdin = nil
	}

	// If the output code writes to stdout.
	if testOutputStreamFile == "" {
		fl, err := os.Create(testOutputFile)
		if err != nil {
			panic(err)
		}
		defer fl.Close()

		testStdout = fl
	} else {
		err := os.Symlink(testOutputFile, filepath.Join(sandboxDir, testOutputStreamFile))
		if err != nil {
			panic(err)
		}
		testStdout = nil
	}

	// Run submission against the testcase.
	timeConsumed, memoryConsumed, err := Execute(sandboxDir, runScript,
		testStdin, testStdout, &testStderr,
		timeLimit, memoryLimit)

	verd.timeConsumed = timeConsumed
	verd.memoryConsumed = memoryConsumed
	verd.stderrLog = testStderr.String()

	if err == context.DeadlineExceeded {
		verd.verdict = color.YellowString("TLE")
	} else if memoryConsumed >= memoryLimit {
		verd.verdict = color.YellowString("MLE")
	} else if err != nil {
		verd.verdict = color.RedString("RTE")

	} else {
		// Program exited gracefully (exit code 0).
		// Run checker to determine verdict.

		checkerScript, err := utils.CleanTemplate(checkerScript, map[string]string{
			"inputFile":    testInputFile,
			"outputFile":   testOutputFile,
			"expectedFile": testExpectedFile,
		})
		if err != nil {
			panic(err)
		}

		var checkerStderr strings.Builder
		_, _, err = Execute(sandboxDir, checkerScript,
			nil, nil, &checkerStderr,
			time.Second*10, 256*1024*1024)

		verd.checkerLog = checkerStderr.String()

		// Determine checker verdict by exit code.
		// Refer github.com/MikeMirzayanov/testlib/blob/master/testlib.h#L219
		// for complete details of all EXIT codes.

		var exitErrno int
		if exitErr, ok := err.(*exec.ExitError); !ok {
			exitErrno = 0
		} else {
			exitErrno = exitErr.ExitCode()
		}

		switch exitErrno {
		case 0, 7:
			// ==> AC
			// ==> POINTS
			verd.verdict = color.GreenString("AC")
		case 1, 2:
			// ==> WA
			// ==> PE
			verd.verdict = color.RedString("WA")
		case 3:
			// ==> FAIL
			verd.verdict = color.HiRedString("FAIL")
		case 8:
			// ==> EOF
			verd.verdict = color.YellowString("EOF")
		}
	}

	// Read input/output/expected files; Truncates to 80 bytes.
	truncReadFile := func(filePath string) string {
		fl, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		buf := make([]byte, 80)
		if n, _ := fl.Read(buf); n == 80 {
			// Output has been truncated.
			buf = append(buf, []byte("...")...)
		}
		return string(bytes.Trim(buf, "\x00"))
	}

	verd.testDetails.truncInput = truncReadFile(testInputFile)
	verd.testDetails.truncOutput = truncReadFile(testOutputFile)
	verd.testDetails.truncExpected = truncReadFile(testExpectedFile)

	return
}

func customTestingMode(sandboxDir, runScript string,
	timeLimit time.Duration, memoryLimit uint64) {

	fmt.Println(color.GreenString("────────────────────────────────"))
	timeConsumed, memoryConsumed, _ := Execute(sandboxDir, runScript,
		os.Stdin, os.Stdout, os.Stderr,
		timeLimit, memoryLimit)
	fmt.Println(color.GreenString("────────────────────────────────"))

	// Time: 32ms    Memory: 2000kb
	c := color.New(color.FgBlue, color.Bold).SprintFunc()
	timeConsumed = timeConsumed.Truncate(time.Millisecond)
	memoryKB := fmt.Sprintf("%dkb", memoryConsumed/1024)
	fmt.Println(c("Time:"), timeConsumed, "\t"+c("Memory:"), memoryKB)
	fmt.Println()
}
