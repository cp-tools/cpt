// +build !windows

package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/kballard/go-shellquote"
	"golang.org/x/sys/unix"
)

// Execute runs the command with resource restrictions.
func Execute(dir, command string,
	stdin io.Reader, stdout, stderr io.Writer,
	timeLimit time.Duration, memoryLimit uint64) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, os.Args[0], "_sandbox",
		"--memory-limit", fmt.Sprint(memoryLimit),
		"--command", fmt.Sprintf("\"%v\"", command))
	cmd.Dir = dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = stdin, stdout, stderr

	return cmd.Run()
}

// Sandbox runs child process with specified memory.
func Sandbox(command string, memoryLimit uint64) {
	// Set memory limit for process.
	var rLimit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_AS, &rLimit); err != nil {
		fmt.Fprintln(os.Stderr, "error getting rlimit:", err)
		os.Exit(1)
	}

	if memoryLimit > rLimit.Max {
		memoryLimit = rLimit.Max
	}
	rLimit.Cur = memoryLimit

	if err := unix.Setrlimit(unix.RLIMIT_AS, &rLimit); err != nil {
		fmt.Fprintln(os.Stderr, "error setting rlimit:", err)
		os.Exit(1)
	}

	cmds, err := shellquote.Split(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
