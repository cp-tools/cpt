package util

import (
	"fmt"
	"reflect"
	"time"

	"github.com/cp-tools/cpt/packages/conf"

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
	fmt.Fprintln(writer, "")
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
