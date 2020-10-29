package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func getLatestReleaseInfo(latestReleaseLink string) (latestVersion, descMsg string) {
	fmt.Println(color.BlueString("Fetching latest release details from GitHub..."))
	// Get information of latest release from github.
	resp, err := http.Get(latestReleaseLink)
	if err != nil {
		fmt.Println(color.RedString("error while querying releases:"), err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(color.RedString("unexpected error occurred:"), err)
		os.Exit(1)
	}

	dataMap := make(map[string]interface{})
	json.Unmarshal(data, &dataMap)

	latestVersion = dataMap["tag_name"].(string)
	descMsg = dataMap["body"].(string)
	return
}

func getReleaseTarball(releaseTarballLink string) *tar.Reader {
	fmt.Println(color.BlueString("Downloading latest release from GitHub..."))
	resp, err := http.Get(releaseTarballLink)
	if err != nil {
		fmt.Println(color.RedString("error while downloading release:"), err)
		os.Exit(1)
	}
	// resp.Body.Close()

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println(color.RedString("unexpected error occurred:"), err)
		os.Exit(1)
	}
	defer gzr.Close()

	trFile := tar.NewReader(gzr)
	if _, err := trFile.Next(); err != nil {
		fmt.Println(color.RedString("error while extracting tarball:"), err)
		os.Exit(1)
	}

	return trFile
}
