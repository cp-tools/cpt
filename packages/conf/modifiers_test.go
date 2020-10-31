package conf_test

import (
	"testing"

	"github.com/cp-tools/cpt/packages/conf"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cnf := conf.New("conf")
	cnf.Load(testMap1)

	assert.Equal(t, map[string]interface{}{
		"top": true,
		"this": map[string]interface{}{
			"is":   float64(984),
			"data": []interface{}{"dog", "cat", "bat"},
		},
		"empty": nil,
	}, cnf.GetAll())
}
func TestSet(t *testing.T) {
	cnf := conf.New("conf")

	cnf.Set("key.val", 35)
	assert.Equal(t, 35, cnf.GetInt("key.val"))

	cnf.Set("key", true)
	assert.Equal(t, true, cnf.GetBool("key"))
	assert.Equal(t, nil, cnf.Get("key.val"))
}

func TestDelete(t *testing.T) {
	cnf := conf.New("conf")

	cnf.Set("key.val", 35)
	cnf.Set("key.test", "golang")
	assert.Equal(t, 35, cnf.GetInt("key.val"))

	cnf.Delete("key.val")
	assert.Equal(t, false, cnf.Has("key.val"))
	assert.Equal(t, true, cnf.Has("key"))

	cnf.Delete("key.test")
	assert.Equal(t, false, cnf.Has("key.test"))
	assert.Equal(t, false, cnf.Has("key"))
}
