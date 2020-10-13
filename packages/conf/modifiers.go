package conf

import (
	"github.com/knadh/koanf/providers/confmap"
)

// Load merges given map into configuration data.
func (cnf *Conf) Load(val map[string]interface{}) {
	cnf.ko.Load(confmap.Provider(val, "."), nil)
}

// LoadDefault merges given map into default data.
func (cnf *Conf) LoadDefault(val map[string]interface{}) {
	cnf.koDefault.Load(confmap.Provider(val, "."), nil)
}

// Set updates the value at given key to val.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) Set(key string, val interface{}) {
	// Create unflattened map from key.
	dataMap := map[string]interface{}{key: val}
	cnf.ko.Load(confmap.Provider(dataMap, "."), nil)
}

// SetDefault updates the default value at given key to val.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) SetDefault(key string, val interface{}) {
	// Create unflattened map from key.
	dataMap := map[string]interface{}{key: val}
	cnf.koDefault.Load(confmap.Provider(dataMap, "."), nil)
}

// Delete deletes the given key from configuration module.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) Delete(key string) {
	// Erase key from configuration.
	cnf.ko.Delete(key)
}

// DeleteDefault deletes the given key from the default configuration module.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) DeleteDefault(key string) {
	// Erase key from default map.
	cnf.koDefault.Delete(key)
}
