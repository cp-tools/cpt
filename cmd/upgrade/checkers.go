package upgrade

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
)

// Checkers upgrades the default checker files.
func Checkers(checkerDir string, cnf *conf.Conf) {
	// Version of current checkers present.
	currentVersion := cnf.GetParent("checker").GetString("version")

	latestVersion, descMsg := getLatestReleaseInfo("https://api.github.com/repos/cp-tools/cpt-checker/releases/latest")
	// Check if current version is outdated.
	if semver.Compare(currentVersion, latestVersion) >= 0 {
		fmt.Println(color.YellowString("(Current version)"), currentVersion, ">=", latestVersion, color.YellowString("(latest version)"))
		return
	}

	fmt.Println(color.GreenString("New version"), latestVersion, color.GreenString("found!"))
	fmt.Println(descMsg)

	var confirm bool
	err := survey.AskOne(&survey.Confirm{
		Message: "Do you wish to upgrade checkers to '" + latestVersion + "'?",
		Default: true,
	}, &confirm)
	util.SurveyOnInterrupt(err)

	if confirm == false {
		return
	}

	// Download release tarball from GitHub.
	releaseTarballLink := fmt.Sprintf("https://github.com/cp-tools/cpt-checker/releases/download/%v/cpt-checker_%v.tar.gz", latestVersion, runtime.GOOS)
	trRdr := getReleaseTarball(releaseTarballLink)

	checkerMap := make(map[string]string)
	tmpCnf := conf.New("newChecker")
	for true {
		hdr, err := trRdr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(color.RedString("error while extracting tarball:"), err)
			os.Exit(1)
		}

		buf, err := ioutil.ReadAll(trRdr)
		if err != nil {
			fmt.Println(color.RedString("unexpected error occurred:"), err)
			os.Exit(1)
		}

		fileName := filepath.Clean(hdr.Name)
		if fileName == "meta.yaml" {
			// meta file of configurations.
			dataMap := make(map[string]interface{})
			if err := yaml.Unmarshal(buf, &dataMap); err != nil {
				fmt.Println(color.RedString("unexpected error occurred:"), err)
				os.Exit(1)
			}

			tmpCnf.Load(dataMap)
		} else if fileName != "." {
			// checker file (executable).
			checkerFile := filepath.Join(checkerDir, fileName)
			file, err := os.Create(checkerFile)
			if err != nil {
				fmt.Println(color.RedString("error while saving checker executable:"), err)
				os.Exit(1)
			}
			defer file.Close()

			if _, err := file.Write(buf); err != nil {
				fmt.Println(color.RedString("unexpected error occurred:"), err)
				os.Exit(1)
			}

			// Set file permissions to make it executable.
			if err := file.Chmod(os.ModePerm); err != nil {
				fmt.Println(color.RedString("unexpected error occurred:"), err)
				os.Exit(1)
			}

			checkerMap[strings.TrimSuffix(fileName, filepath.Ext(fileName))] = checkerFile
			fmt.Println(color.GreenString("Checker"), fileName, color.GreenString("saved successfully!"))
		}
	}

	for k, v := range checkerMap {
		script := tmpCnf.GetString("checker.checkers." + k + ".script")
		script = strings.ReplaceAll(script, "$FILEPATH", v)
		tmpCnf.Set("checker.checkers."+k+".script", script)
	}

	cnf.GetParent("checker").Load(tmpCnf.GetAll()).WriteFile()
}
