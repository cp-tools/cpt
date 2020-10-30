package upgrade

import (
	"fmt"
	"os"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/inconshreveable/go-update"
	"golang.org/x/mod/semver"
)

// Self upgrades the cli tool.
func Self(currentVersion string) {
	latestVersion, descMsg := getLatestReleaseInfo("https://api.github.com/repos/cp-tools/cpt/releases/latest")
	// Check if current version is outdated.
	if semver.Compare(currentVersion, latestVersion) >= 0 {
		fmt.Println(color.YellowString("(Current version)"), currentVersion, ">=", latestVersion, color.YellowString("(latest version)"))
		return
	}

	fmt.Println(color.GreenString("New version"), latestVersion, color.GreenString("found!"))
	fmt.Println(descMsg)

	var confirm bool
	survey.AskOne(&survey.Confirm{
		Message: "Do you wish to upgrade to '" + latestVersion + "'?",
		Default: true,
	}, &confirm)

	if confirm == false {
		return
	}

	// Download release tarball from GitHub.
	releaseTarballLink := fmt.Sprintf("https://github.com/cp-tools/cpt/releases/download/%v/cpt_%v_%v.tar.gz", latestVersion, runtime.GOOS, runtime.GOARCH)
	trRdr := getReleaseTarball(releaseTarballLink)
	if _, err := trRdr.Next(); err != nil {
		fmt.Println(color.RedString("error while extracting tarball:"), err)
		os.Exit(1)
	}
	// Tarball MUST contain exactly 1 file, the executable.

	// Overwrite current binary with downloaded binary.
	if err := update.Apply(trRdr, update.Options{}); err != nil {
		fmt.Println(color.RedString("error while upgrading to latest version:"), err)
		// Failed to update binary. Rollback to current version.
		if err := update.RollbackError(err); err != nil {
			// This is fatal. Should never happen.
			fmt.Println(color.RedString("error while rolling back to previous version:"), err)
			os.Exit(1)
		}

		fmt.Println(color.YellowString("Rolled back to previous version."))
		return
	}

	fmt.Println(color.GreenString("Successfully upgraded to"), latestVersion, "!")
}
