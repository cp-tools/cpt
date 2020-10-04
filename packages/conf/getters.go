package conf

// Get returns interface{} value of a given key path,
// or nil if key does not exist or is invalid.
func (cnf *Conf) Get(key string) interface{} {
	return cnf.ko.Get(key)
}

// GetInt returns int value of a given key path,
// or 0 if key does not exist or is invalid.
func (cnf *Conf) GetInt(key string) int {
	return cnf.ko.Int(key)
}

// GetString returns string value of a given key path,
// or "" if key does not exist or is invalid.
func (cnf *Conf) GetString(key string) string {
	return cnf.ko.String(key)
}

// GetBool returns bool value of a given key path,
// or false if key does not exist or is invalid.
func (cnf *Conf) GetBool(key string) bool {
	return cnf.ko.Bool(key)
}

// GetFloat64 returns float64 value of a given key path,
// or 0 if key does not exist or is invalid.
func (cnf *Conf) GetFloat64(key string) float64 {
	return cnf.ko.Float64(key)
}
