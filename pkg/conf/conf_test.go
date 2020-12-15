package conf_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var testMap1 = map[string]interface{}{
	"top": true,
	"this": map[string]interface{}{
		"is":   984,
		"data": []string{"dog", "cat", "bat"},
	},
	"empty": nil,
}

var testMap2 = map[string]interface{}{
	"top": false,
	"this": map[string]interface{}{
		"key":  58,
		"data": []int{1, 2, 3},
	},
	"bottom": true,
}

var testMap3 = map[string]interface{}{
	"int":     420,
	"bool":    true,
	"float":   3.14,
	"string":  "where banana",
	"strings": []string{"cat", "dog"},
	"map": map[string]interface{}{
		"int":     910,
		"bool":    false,
		"float":   2.71,
		"string":  "banana there",
		"strings": []string{"zebra", "goat"},
		"map": map[string]interface{}{
			"exists": true,
		},
	},
}

func TestSetParent(t *testing.T) {
	cnf := conf.New("root")
	cnf.Load(testMap1)
	cnf = conf.New("child").SetParent(cnf)
	cnf.Load(testMap2)

	assert.EqualValues(t, (*conf.Conf)(nil), cnf.GetParent("invalid"))
	assert.Panics(t, func() { conf.New("child").SetParent(cnf) })

	assert.EqualValues(t, map[string]interface{}{
		"top": true,
		"this": map[string]interface{}{
			"is":   float64(984),
			"data": []interface{}{"dog", "cat", "bat"},
		},
		"empty": nil,
	}, cnf.GetParent("root").GetAll())
}

func TestLoadFile(t *testing.T) {
	out, _ := yaml.Marshal(testMap1)

	file, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(file.Name())

	file.Write(out)

	cnf := conf.New("conf").LoadFile(file.Name())

	assert.EqualValues(t, nil, cnf.Get("this.key"))
	assert.EqualValues(t, []string{"dog", "cat", "bat"}, cnf.GetStrings("this.data"))
}

func TestWriteFile(t *testing.T) {
	file, _ := ioutil.TempFile(os.TempDir(), "")
	defer os.Remove(file.Name())

	cnf := conf.New("conf").Load(testMap2)
	cnf.LoadFile(file.Name()).WriteFile()

	buf, _ := ioutil.ReadFile(file.Name())
	testMapBuf, _ := yaml.Marshal(testMap2)

	assert.EqualValues(t, testMapBuf, buf)
}
