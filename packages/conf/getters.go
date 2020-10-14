package conf

import "github.com/knadh/koanf/maps"

// Get returns interface{} value of a given key path,
// or nil if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the default data is searched for the same.
func (cnf *Conf) Get(key string) interface{} {
	if !cnf.ko.Exists(key) {
		return cnf.koDefault.Get(key)
	}
	return cnf.ko.Get(key)
}

// GetAll merges the configured values with the default
// values and returns the data as a map.
func (cnf *Conf) GetAll() map[string]interface{} {
	mp := cnf.ko.Raw()
	maps.Merge(cnf.koDefault.Raw(), mp)
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
	return cnf.ko.MapKeys(key)
}
