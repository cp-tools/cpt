package conf

import (
	"github.com/knadh/koanf/providers/confmap"
)

// Load merges given map into configuration data.
func (cnf *Conf) Load(val map[string]interface{}) {
	cnf.ko.Load(confmap.Provider(val, "."), nil)
}

// Set updates the value at given key to val.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) Set(key string, val interface{}) {
	// Create unflattened map from key.
	dataMap := map[string]interface{}{key: val}
	cnf.ko.Load(confmap.Provider(dataMap, "."), nil)
}

// Erase deletes the given key from configuration module.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) Erase(key string) {
	// Erase key from configuration.
	cnf.ko.Erase(key)
}
