package conf

import (
	"github.com/knadh/koanf/providers/confmap"
)

// Load merges given map into configuration data.
func (cnf *Conf) Load(val map[string]interface{}) *Conf {
	cnf.ko.Load(confmap.Provider(val, "."), nil)
	return cnf
}

// LoadDefault merges given map into default data.
func (cnf *Conf) LoadDefault(val map[string]interface{}) *Conf {
	cnf.koDefault.Load(confmap.Provider(val, "."), nil)
	return cnf
}

// Set updates the value at given key to val.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) Set(key string, val interface{}) *Conf {
	// Create unflattened map from key.
	dataMap := map[string]interface{}{key: val}
	cnf.ko.Load(confmap.Provider(dataMap, "."), nil)
	return cnf
}

// SetDefault updates the default value at given key to val.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) SetDefault(key string, val interface{}) *Conf {
	// Create unflattened map from key.
	dataMap := map[string]interface{}{key: val}
	cnf.koDefault.Load(confmap.Provider(dataMap, "."), nil)
	return cnf
}

// Delete deletes the given key from configuration module.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) Delete(key string) *Conf {
	// Erase key from configuration.
	cnf.ko.Delete(key)
	return cnf
}

// DeleteDefault deletes the given key from the default configuration module.
//
// The key should be a flattened path, with '.' as the delim.
func (cnf *Conf) DeleteDefault(key string) *Conf {
	// Erase key from default map.
	cnf.koDefault.Delete(key)
	return cnf
}
