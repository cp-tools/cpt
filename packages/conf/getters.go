package conf

import (
	"github.com/knadh/koanf/maps"
)

// GetParent returns parent of conf with given name.
// Returns nil if no parent of given name exists.
func (cnf *Conf) GetParent(name string) *Conf {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.name == name {
			// Found it. Break.
			break
		}

		rootCnf = rootCnf.parentCnf
	}

	return rootCnf
}

// Get returns interface{} value of a given key path,
// or nil if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) Get(key string) interface{} {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			// Found conf with given key.
			break
		}

		rootCnf = rootCnf.parentCnf
	}

	if rootCnf == nil {
		// key does not exist.
		return nil
	}

	// Recursively merge data if is map.
	if mp, ok := rootCnf.ko.Get(key).(map[string]interface{}); ok {
		rootCnf = rootCnf.parentCnf
		for rootCnf != nil {
			if mp1, ok := rootCnf.ko.Get(key).(map[string]interface{}); ok {
				maps.Merge(mp, mp1)
				mp = mp1
			}

			rootCnf = rootCnf.parentCnf
		}

		return mp
	}
	// Not map, return value.
	return rootCnf.ko.Get(key)
}

// Has returns true if the given key exists in configuration.
func (cnf *Conf) Has(key string) bool {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			return true
		}

		rootCnf = rootCnf.parentCnf
	}

	return false
}

// GetAll merges the configured values with the default
// values and returns the data as a map.
func (cnf *Conf) GetAll() map[string]interface{} {
	mp := make(map[string]interface{})

	rootCnf := cnf
	for rootCnf != nil {
		mp1 := rootCnf.ko.Raw()
		maps.Merge(mp, mp1)
		mp = mp1

		rootCnf = rootCnf.parentCnf
	}

	return mp
}

// GetInt returns int value of a given key path,
// or 0 if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) GetInt(key string) int {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			return rootCnf.ko.Int(key)
		}

		rootCnf = rootCnf.parentCnf
	}

	return 0
}

// GetString returns string value of a given key path,
// or "" if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) GetString(key string) string {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			return rootCnf.ko.String(key)
		}

		rootCnf = rootCnf.parentCnf
	}

	return ""
}

// GetStrings returns []string slice value of a given key path,
// or "" if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) GetStrings(key string) []string {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			return rootCnf.ko.Strings(key)
		}

		rootCnf = rootCnf.parentCnf
	}

	return []string{}
}

// GetBool returns bool value of a given key path,
// or false if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) GetBool(key string) bool {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			return rootCnf.ko.Bool(key)
		}

		rootCnf = rootCnf.parentCnf
	}

	return false
}

// GetFloat64 returns float64 value of a given key path,
// or 0 if key does not exist or is invalid.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) GetFloat64(key string) float64 {
	rootCnf := cnf
	for rootCnf != nil {
		if rootCnf.ko.Exists(key) {
			return rootCnf.ko.Float64(key)
		}

		rootCnf = rootCnf.parentCnf
	}

	return 0
}

// GetMapKeys returns a string list of keys in a map addressed
// by the given path. If the path is not a map, an empty
// string slice is returned.
//
// If given key doesn't exist in configuration module,
// the chained parent data is searched for the same.
func (cnf *Conf) GetMapKeys(key string) []string {
	data := make([]string, 0)

	if mp, ok := cnf.Get(key).(map[string]interface{}); ok {
		for key := range mp {
			data = append(data, key)
		}
	}

	return data
}
