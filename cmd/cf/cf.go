package cf

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cp-tools/cpt-lib/codeforces"

	"github.com/fatih/color"
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

// some ugly functions exclusive to codeforces, below!
func colorVerdict(verdict string) string {
	verdict = strings.ReplaceAll(verdict, "Time limit exceeded", "TLE")
	verdict = strings.ReplaceAll(verdict, "Compilation error", "CE")
	verdict = strings.ReplaceAll(verdict, "Runtime error", "RTE")
	verdict = strings.ReplaceAll(verdict, "Wrong answer", "WA")
	verdict = strings.ReplaceAll(verdict, "Accepted", "AC")
	verdict = strings.ReplaceAll(verdict, "Partial result", "PR")

	if strings.HasPrefix(verdict, "CE") || strings.HasPrefix(verdict, "TLE") {
		verdict = color.HiYellowString(verdict)
	} else if strings.HasPrefix(verdict, "RTE") || strings.HasPrefix(verdict, "WA") {
		verdict = color.HiRedString(verdict)
	} else if verdict == "AC" || strings.HasPrefix(verdict, "PR") {
		verdict = color.HiGreenString(verdict)
	}
	return verdict
}
