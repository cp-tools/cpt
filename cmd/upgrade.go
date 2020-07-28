package cmd

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
	Short: "Upgrade binary using github releases",
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
	fmt.Println("Querying releases information from github")

	checkerFlag, _ := lflags.GetBool("checker")
	if checkerFlag == true {
		// get checker releases from github
		const releasesURL = "https://api.github.com/repos/cp-tools/cpt-checker/releases"
		allReleases := getReleases(releasesURL)

		fmt.Println("Current checker build version:", viper.GetString("checker_version"))

		var upgVers string
		util.SurveyErr(survey.AskOne(&survey.Select{
			Message: "Select version to upgrade/degrade to:",
			Options: allReleases,
		}, &upgVers))

		fmt.Println("Downloading checkers. Please wait...")
		fileURL := fmt.Sprintf("https://github.com/cp-tools/cpt-checker/releases/download/%v/cpt-checker_%v.tar.gz",
			upgVers, runtime.GOOS)
		resp, err := http.Get(fileURL)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		tr := getTarball(resp)

		fmt.Println("Saved checkers:")
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
			fmt.Println("Failed to save configurations")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Checkers upgraded successfully!")
	} else {
		// get cpt releases from github
		const releasesURL = "https://api.github.com/repos/cp-tools/cpt/releases"
		allReleases := getReleases(releasesURL)

		fmt.Println("Current cli app version:", rootCmd.Version)

		var upgVers string
		util.SurveyErr(survey.AskOne(&survey.Select{
			Message: "Select version to upgrade/degrade to:",
			Options: allReleases,
		}, &upgVers))

		fmt.Println("Downloading cli binary. Please wait...")
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
				fmt.Println("Failed to roll back from defective update")
				fmt.Println(rerr)
				os.Exit(1)
			}

			fmt.Println("Could not update binary executable")
			fmt.Println("Rolled back to previous version")
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("Upgraded cli app successfully!")
	}
}

func getReleases(releasesURL string) []string {
	// fetch data of all releases
	resp, err := http.Get(releasesURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var allReleases []string
	for _, tag := range gjson.GetBytes(body, "#.tag_name").Array() {
		allReleases = append(allReleases, tag.String())
	}
	return allReleases
}

func getTarball(resp *http.Response) *tar.Reader {
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println("Could not read tarball")
		fmt.Println(err)
		os.Exit(1)
	}
	defer gzr.Close()
	return tar.NewReader(gzr)
}
