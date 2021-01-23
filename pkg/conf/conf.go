package conf

import (
	"log"
	"os"

	"github.com/cp-tools/cpt/utils"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Conf is the configuration module.
type Conf struct {
	name   string
	ko     *koanf.Koanf
	koFile string

	parentCnf *Conf
}

// New returns a new instance of Conf.
func New(name string) *Conf {
	cnf := new(Conf)
	cnf.name = name
	cnf.ko = koanf.New(".")

	return cnf
}

// SetParent sets the parent conf.
func (cnf *Conf) SetParent(parentCnf *Conf) *Conf {
	cnf.parentCnf = parentCnf
	// Check for conf name clash.
	rootCnf := cnf.parentCnf
	for rootCnf != nil {
		if cnf.name == rootCnf.name {
			panic("name clash occurred")
		}
		rootCnf = rootCnf.parentCnf
	}

	return cnf
}

// LoadFile reads and loads data from file at
// given path to the configuration module.
// Does nothing if file doesn't exist.
//
// Ensure the file at given path is of YAML format.
func (cnf *Conf) LoadFile(path string) *Conf {
	cnf.koFile = path
	// Check if file at given path exists.
	if !utils.FileExists(path) {
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
	file, err := os.Create(cnf.koFile)
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
