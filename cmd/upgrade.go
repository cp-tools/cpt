package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"golang.org/x/mod/semver"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade binary using github releases",
	Run: func(cmd *cobra.Command, args []string) {
		upgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() {
	var allReleases []string
	fmt.Println("Querying releases information from github")
	// fetch data of all releases
	const releasesURL = "https://api.github.com/repos/cp-tools/cpt/releases"
	resp, err := http.Get(releasesURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, tag := range gjson.GetBytes(body, "#.tag_name").Array() {
		allReleases = append(allReleases, tag.String())
	}

	if len(allReleases) == 0 || semver.Compare(rootCmd.Version, allReleases[0]) == 0 {
		fmt.Println("Current version", rootCmd.Version, "is the latest!")

		var choice bool
		util.SurveyErr(survey.AskOne(&survey.Confirm{
			Message: "Proceed to version selection menu?",
		}, &choice))
		if choice == false {
			os.Exit(0)
		}
	}

	// give user list of options to select from
	fmt.Println("Current version:", rootCmd.Version)

	var upgVers string
	util.SurveyErr(survey.AskOne(&survey.Select{
		Message: "Select version you wish to upgrade/degrade to:",
		Options: allReleases,
	}, &upgVers))

	// fetch binary zip file
	fmt.Println("Downloading binary file. Please wait...")

	binaryURL := fmt.Sprintf("https://github.com/cp-tools/cpt/releases/download/%v/cpt_%v_%v.tar.gz",
		upgVers, runtime.GOOS, runtime.GOARCH)

	resp, err = http.Get(binaryURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// extract binary app from zip file
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println("Could not read tarball")
		fmt.Println(err)
		os.Exit(1)
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	tr.Next()

	// update binary app
	err = update.Apply(tr, update.Options{})
	if err != nil {
		// update failed. Attempt rollback
		if rerr := update.RollbackError(err); rerr != nil {
			fmt.Println("Failed to roll back from defective update")
			fmt.Println(rerr)
			os.Exit(1)
		}

		fmt.Println("Could not update binary executable")
		fmt.Println("Rolled back to previous version")
		fmt.Println(err)
	} else {
		fmt.Println("Successfully updated to", upgVers)
	}
}
