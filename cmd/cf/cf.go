package cf

import (
	"net/http"
	"path/filepath"

	"github.com/cp-tools/cpt-lib/codeforces"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// RootCmd is 'cf' subcommand
	RootCmd = &cobra.Command{
		Use:     "cf",
		Aliases: []string{"codeforces"},
		Short:   "Utilities common to codeforces",
	}

	cfViper = viper.New()
)

// InitConfig loads configurations
func InitConfig(cfgDir string) {
	// load global settings
	cfgFile := filepath.Join(cfgDir, "cf.json")
	cfViper.SetConfigFile(cfgFile)
	cfViper.SafeWriteConfig()
	cfViper.ReadInConfig()

	codeforces.SessCln = http.DefaultClient
}
