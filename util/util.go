package util

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"

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
