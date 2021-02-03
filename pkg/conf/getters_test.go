package conf_test

import (
	"testing"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/stretchr/testify/assert"
)

var globalCnf = conf.New("conf").Load(testMap3)

func TestGetParent(t *testing.T) {
	cnf1 := conf.New("c1")
	cnf2 := conf.New("c2")
	cnf3 := conf.New("c3")
	cnf4 := conf.New("c4")

	cnf2.SetParent(cnf1)
	cnf3.SetParent(cnf2)
	cnf4.SetParent(cnf2)

	assert.Equal(t, cnf1, cnf2.GetParent("c1"))
	assert.Equal(t, cnf2, cnf3.GetParent("c2"))
	assert.Equal(t, cnf2, cnf4.GetParent("c2"))
	assert.Equal(t, (*conf.Conf)(nil), cnf1.GetParent("invalid"))
}

func TestGet(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, float64(910), cnf.Get("map.int"))
	assert.Equal(t, nil, cnf.Get("invalid.key"))

	cnf.Set("map.float", 1729)
	assert.Equal(t, float64(1729), cnf.Get("map.float"))

	cnf.Load(map[string]interface{}{
		"map": map[string]interface{}{
			"strings": []string{"pig", "horse"},
		},
		"empty": nil,
	})

	assert.Equal(t, map[string]interface{}{
		"int":     float64(420),
		"bool":    true,
		"float":   float64(3.14),
		"string":  "where banana",
		"strings": []interface{}{"cat", "dog"},
		"map": map[string]interface{}{
			"int":     float64(910),
			"bool":    false,
			"float":   float64(1729),
			"string":  "banana there",
			"strings": []interface{}{"pig", "horse"},
			"map": map[string]interface{}{
				"exists": true,
			},
		},
		"empty": nil,
	}, cnf.Get(""))
}

func TestHas(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, true, cnf.Has("map.map.exists"))
	assert.Equal(t, false, cnf.Has("map.exists"))
	assert.Equal(t, false, cnf.Has("invalid.key"))
}

func TestGetInt(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, int(910), cnf.GetInt("map.int"))
	assert.Equal(t, int(0), cnf.GetInt("map.string"))
	assert.Equal(t, int(0), cnf.GetInt("invalid.key"))
}

func TestGetString(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, "where banana", cnf.GetString("string"))
	assert.Equal(t, "[cat dog]", cnf.GetString("strings"))
	assert.Equal(t, "true", cnf.GetString("map.map.exists"))
	assert.Equal(t, "", cnf.GetString("invalid.key"))
}

func TestGetStrings(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, []string{}, cnf.GetStrings("string"))
	assert.Equal(t, []string{"zebra", "goat"}, cnf.GetStrings("map.strings"))
	assert.Equal(t, []string{}, cnf.GetStrings("invalid.key"))
}

func TestGetBool(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, true, cnf.GetBool("bool"))
	assert.Equal(t, false, cnf.GetBool("map.bool"))
	assert.Equal(t, false, cnf.GetBool("int"))
	assert.Equal(t, false, cnf.GetBool("invalid.key"))
}

func TestGetFloat64(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, float64(3.14), cnf.GetFloat64("float"))
	assert.Equal(t, float64(2.71), cnf.GetFloat64("map.float"))
	assert.Equal(t, float64(0), cnf.GetFloat64("string"))
	assert.Equal(t, float64(0), cnf.GetFloat64("invalid.key"))
}

func TestGetMapKeys(t *testing.T) {
	cnf := conf.New("cnf").SetParent(globalCnf)

	assert.Equal(t, []string{
		"bool", "float", "int", "map", "string", "strings",
	}, cnf.GetMapKeys(""))
	assert.Equal(t, []string{"exists"}, cnf.GetMapKeys("map.map"))
	assert.Equal(t, []string{}, cnf.GetMapKeys("float"))
}
