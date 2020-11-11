package util

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/cp-tools/cpt/packages/conf"

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

// LoadLocalConf returns local folder conf.
func LoadLocalConf(cnf *conf.Conf) *conf.Conf {
	cnf = conf.New("local").SetParent(cnf)
	cnf.LoadFile("meta.yaml")

	return cnf
}

// ColorHeaderFormat sets color (for headers mostly).
func ColorHeaderFormat(str ...string) []string {
	for i := range str {
		str[i] = color.New(color.FgBlue, color.Bold, color.Underline).Sprint(str[i])
	}
	return str
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
