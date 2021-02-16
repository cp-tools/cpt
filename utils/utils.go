package utils

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

// RunCountdown runs countdown with a static message.
func RunCountdown(dur time.Duration, msg string) {
	writer := uilive.New()
	writer.Start()
	for ; dur.Seconds() > 0; dur -= time.Second {
		fmt.Fprintln(writer, msg, dur.String())
		time.Sleep(time.Second)
	}
	fmt.Fprint(writer)
	writer.Stop()
}

// ExtractMapKeys returns top-level keys of map.
func ExtractMapKeys(varMap interface{}) (data []string) {
	keys := reflect.ValueOf(varMap).MapKeys()
	for _, key := range keys {
		data = append(data, key.String())
	}
	return
}

// CleanTemplate creates and runs template on passed string, with given params.
func CleanTemplate(str string, data interface{}) (string, error) {
	tmplt, err := template.New("").Parse(str)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	if err := tmplt.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}

// SurveyOnInterrupt is run on SIGINT.
func SurveyOnInterrupt(err error) {
	if err == terminal.InterruptErr {
		fmt.Println("interrupted")
		os.Exit(130)
	} else if err != nil {
		fmt.Println(color.RedString("unexpected error occurred:"), err)
		os.Exit(1)
	}
}

// FileExists returns a bool signifying if given file exists.
func FileExists(filename string) bool {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !f.IsDir()
}

// OpenURL opens the specified URL in the default browser of the user.
// Attribution: https://stackoverflow.com/a/39324149
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
