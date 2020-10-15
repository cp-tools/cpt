package conf

import (
	"sort"

	"github.com/knadh/koanf/maps"
)

// Get returns interface{} value of a given key path,
// or nil if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) Get(key string) interface{} {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.Get(key)
	}
	// Merge default and configurations if key is map.
	if mp1, ok := cnf.koDefault.Get(key).(map[string]interface{}); ok {
		if mp2, ok := cnf.ko.Get(key).(map[string]interface{}); ok {
			maps.Merge(mp2, mp1)
			return mp1
		}
	}

	return cnf.ko.Get(key)
}

// Has returns true if the given key exists in configuration.
func (cnf *Conf) Has(key string) bool {
	// Check if key exists in either.
	if cnf.ko.Exists(key) || cnf.koDefault.Exists(key) {
		return true
	}
	return false
}

// GetAll merges the configured values with the default
// values and returns the data as a map.
func (cnf *Conf) GetAll() map[string]interface{} {
	mp := cnf.koDefault.Raw()
	maps.Merge(cnf.ko.Raw(), mp)
	return mp
}

// GetInt returns int value of a given key path,
// or 0 if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) GetInt(key string) int {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.Int(key)
	}
	return cnf.ko.Int(key)
}

// GetString returns string value of a given key path,
// or "" if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) GetString(key string) string {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.String(key)
	}
	return cnf.ko.String(key)
}

// GetStrings returns []string slice value of a given key path,
// or "" if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) GetStrings(key string) []string {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.Strings(key)
	}
	return cnf.ko.Strings(key)
}

// GetBool returns bool value of a given key path,
// or false if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) GetBool(key string) bool {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.Bool(key)
	}
	return cnf.ko.Bool(key)
}

// GetFloat64 returns float64 value of a given key path,
// or 0 if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) GetFloat64(key string) float64 {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.Float64(key)
	}
	return cnf.ko.Float64(key)
}

// GetMapKeys returns a sorted string list of keys in a map
// addressed by the given path. If the path is not a map,
// an empty string slice is returned.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) GetMapKeys(key string) []string {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.MapKeys(key)
	}
	data := cnf.koDefault.MapKeys(key)
	data = append(data, cnf.ko.MapKeys(key)...)
	sort.Strings(data)
	return data
}
