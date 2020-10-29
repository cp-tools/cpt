package cmd

/*
import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cp-tools/cpt/util"
	"github.com/fatih/color"

	"github.com/AlecAivazis/survey/v2"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Args:  cobra.NoArgs,
	Short: "Upgrade binary using github releases",
	Long: `Downloads latest release from github and replaces current executable.
The easiest way to upgrade to the latest release, with rollback option on failure.

Presents you with a selection list of ALL releases so far, giving you the greatest
flexibility to choose the version you want to use. Using latest version is recommended.

Use --checker flag to upgrade the set of default checkers available for 'cpt test'.
Note that, checker upgrade replaces checkers in $CONFIGDIR/cpt/checkers with the same name.

Usage examples:
cpt upgrade
                            Bring selection menu for version to upgrade to
cpt upgrade -c
                            Same as 'cpt upgrade' but for default checkers

`,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// add flags here
	upgradeCmd.Flags().BoolP("checker", "c", false, "Upgrade default checkers")

	upgradeCmd.Run = func(cmd *cobra.Command, args []string) {
		lflags := upgradeCmd.Flags()
		upgrade(lflags)
	}

}

func upgrade(lflags *pflag.FlagSet) {
	color.Blue("Querying releases information from github")

	checkerFlag := lflags.MustGetBool("checker")
	if checkerFlag == true {
		// get checker releases from github
		const releasesURL = "https://api.github.com/repos/cp-tools/cpt-checker/releases"
		allReleases := getReleases(releasesURL)

		fmt.Println(color.BlueString("Current checker build version:"), viper.GetString("checker_version"))

		var upgVers string
		util.SurveyErr(survey.AskOne(&survey.Select{
			Message: "Select version to upgrade/degrade to:",
			Options: allReleases,
		}, &upgVers))

		color.Green("Downloading checkers. Please wait...")
		fileURL := fmt.Sprintf("https://github.com/cp-tools/cpt-checker/releases/download/%v/cpt-checker_%v.tar.gz",
			upgVers, runtime.GOOS)
		resp, err := http.Get(fileURL)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		tr := getTarball(resp)

		color.Blue("Downloaded checkers:")
		for true {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			// we know that everything inside the tar is a file
			fl, err := os.Create(filepath.Join(cfgDir, "checkers", hdr.Name))
			if err != nil {
				panic(err)
			}
			// write file from tar to
			if _, err = io.Copy(fl, tr); err != nil {
				panic(err)
			}
			// change mode permissions
			if err = fl.Chmod(os.ModePerm); err != nil {
				panic(err)
			}
			fl.Close()
			fmt.Printf("  - %v\n", hdr.Name)
		}

		viper.Set("checker_version", upgVers)
		if err := viper.WriteConfig(); err != nil {
			color.Red("Failed to save configurations")
			fmt.Println(err)
			os.Exit(1)
		}

		color.Green("Checkers upgraded successfully!")
	} else {
		// get cpt releases from github
		const releasesURL = "https://api.github.com/repos/cp-tools/cpt/releases"
		allReleases := getReleases(releasesURL)

		fmt.Println(color.BlueString("Current cli app version:"), rootCmd.Version)

		var upgVers string
		util.SurveyErr(survey.AskOne(&survey.Select{
			Message: "Select version to upgrade/degrade to:",
			Options: allReleases,
		}, &upgVers))

		color.Blue("Downloading cli binary. Please wait...")
		fileURL := fmt.Sprintf("https://github.com/cp-tools/cpt/releases/download/%v/cpt_%v_%v.tar.gz",
			upgVers, runtime.GOOS, runtime.GOARCH)
		resp, err := http.Get(fileURL)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		tr := getTarball(resp)
		tr.Next()

		// we know that only 1 file exists inside tar
		if err := update.Apply(tr, update.Options{}); err != nil {
			// update failed. Attempt rollback
			if rerr := update.RollbackError(err); rerr != nil {
				color.Red("Failed to roll back from stalled update")
				fmt.Println(rerr)
				os.Exit(1)
			}

			color.Red("Could not update binary executable")
			color.Yellow("Rolled back to previous version")
			fmt.Println(err)
			os.Exit(0)
		}
		color.Green("Upgraded cli app successfully!")
	}
}
*/
