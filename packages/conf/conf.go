package conf

import (
	"log"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Conf is the configuration module.
type Conf struct {
	ko         *koanf.Koanf
	koDefault  *koanf.Koanf
	koFilePath string
}

// New returns a new instance of Conf.
func New() *Conf {
	cnf := new(Conf)
	cnf.ko = koanf.New(".")
	cnf.koDefault = koanf.New(".")
	return cnf
}

// LoadFile reads and loads data from file at
// given path to the configuration module.
// Does nothing if file doesn't exist.
//
// Ensure the file at given path is of YAML format.
func (cnf *Conf) LoadFile(path string) *Conf {
	cnf.koFilePath = path
	// Check if file at given path exists.
	if file, err := os.Stat(path); os.IsNotExist(err) || file.IsDir() {
		return cnf
	}
	// Load YAML conf file.
	if err := cnf.ko.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Fatalf("error loading conf file: %v", err)
	}
	return cnf
}

// WriteFile overwrites data from the configuration module
// to the file last set using LoadConf().
// Does nothing if configuration data is empty.
// Values from the default map are not written.
//
// The written data is of YAML format.
func (cnf *Conf) WriteFile() *Conf {
	// Raw data of configuration to write.
	rawMap := cnf.ko.Raw()

	if len(rawMap) == 0 {
		return cnf
	}
	// Create file if it does not exist,
	// and truncate the file if it does.
	file, err := os.Create(cnf.koFilePath)
	if err != nil {
		log.Fatalf("error creating conf file: %v", err)
	}
	defer file.Close()

	// Marshal conf to YAML format.
	data, err := yaml.Parser().Marshal(rawMap)
	if err != nil {
		log.Fatalf("unexpected error occurred: %v", err)
	}
	// Write data to conf file.
	if _, err := file.Write(data); err != nil {
		log.Fatalf("error writing to conf file: %v", err)
	}
	return cnf
}
