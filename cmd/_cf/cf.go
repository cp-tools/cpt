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
		Use:   "cf",
		Short: "Utilities common to codeforces",
		Long: `Helper functions for codeforces are in this subcommand.
Configure codeforces login credentials using 'cpt cf config'. Your
password is encrypted with AES and saved to $CONFIGDIR/cpt/cf.json.

However, note that anyone who gets access to the file could extract
your password using the corresponding key. The encryption is to prevent
people from directly reading the file if saved as plaintext.

You will find argument [SPECIFIER] required in multiple places. An imcomplete list
of valid specifiers for this subcommand are as follows:

- Links:
  codeforces.com/contest/1388
  https://codeforces.com/gym/102672
  codeforces.com/gym/102672/problem/I
  codeforces.com/problemset/problem/1389/A
  codeforces.com/group/OzCWQ49fxc/contest/279141

- Direct:
  1234 f
  102672
  OzCWQ49fxc 279141e

If no specifier is provided, specifier is parsed from current
directory structure.
`,
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
