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
)

type (
	testExecDetails struct {
		verdict    string
		failLog    string
		checkerLog string
		stderrLog  string
		runtimeLog string

		timeConsumed   time.Duration
		memoryConsumed uint64

		testDetails struct {
			truncInput    string
			truncOutput   string
			truncExpected string
		}
	}
)

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
		os.Remove(filepath.Join(sandboxDir, testInputStreamFile))
		err := os.Link(testInputFile, filepath.Join(sandboxDir, testInputStreamFile))
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
		os.Remove(filepath.Join(sandboxDir, testOutputStreamFile))
		err := os.Link(testOutputFile, filepath.Join(sandboxDir, testOutputStreamFile))
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
		verd.runtimeLog = err.Error()
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

func interactiveTestingMode(sandboxDir, runScript string,
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
