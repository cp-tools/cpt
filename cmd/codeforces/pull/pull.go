package pull

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/cmd/generate"
	"github.com/cp-tools/cpt/packages/conf"

	"github.com/fatih/color"
)

// Pull pulls
func Pull(arg codeforces.Args, mode, username string, cnf *conf.Conf) {
	// Fetch submission details matching arg.
	fmt.Println(color.BlueString("Fetching submission details of:"), arg)
	chanSubmissions, err := arg.GetSubmissions(username, 100) // Anything more than this is NO NO!
	if err != nil {
		fmt.Println(color.RedString("error occurred while fetching submissions:"), err)
		os.Exit(1)
	}

	// Template to use to specify folder path for each fetched problem submission.
	folderPath := filepath.Join(cnf.GetStrings("pull.problemFolderPath")...)
	folderPathTmplt, err := template.New("template").Parse(folderPath)
	if err != nil {
		fmt.Println(color.RedString("error occurred while parsing 'folderPath' template:"), err)
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(color.RedString("unexpected error occurred:"), err)
		os.Exit(1)
	}

	mp := make(map[codeforces.Args]bool)
	for submissions := range chanSubmissions {
		for _, sub := range submissions {
			if shouldPull(sub, mode, mp) == false {
				continue
			}
			// Fetch submission source code.
			fmt.Println(color.GreenString("pulling source code of:"), sub.Arg, "-", sub.ID)
			sourceCode, err := sub.GetSourceCode()
			if err != nil {
				fmt.Println("error occurred while pulling source code:", err)
				// Don't quit; continue.
				continue
			}

			// Determine folder path to parse problem submission to.
			var problemDirBuf strings.Builder
			folderPathTmplt.Execute(&problemDirBuf, sub)
			problemDir := filepath.Clean(problemDirBuf.String())
			// Create folder and check for errors.
			if err := os.MkdirAll(problemDir, os.ModePerm); err != nil {
				fmt.Println(color.RedString("error occurred while creating problem folder:"), err)
				os.Exit(1)
			}

			// Move into problem directory.
			if err := os.Chdir(problemDir); err != nil {
				fmt.Println(color.RedString("unexpected error occurred:"), err)
				os.Exit(1)
			}

			// Decide baseFileName and fileExtension.
			baseFileName := filepath.Base(problemDir)
			fileExtension := codeforces.LanguageExtn[sub.Language]

			// Determine fileName and write sourceCode to it.
			fileName := generate.DecideFileName(baseFileName, fileExtension)
			file, err := os.Create(fileName)
			if err != nil {
				color.Red("error creating code file: %v", err)
				os.Exit(1)
			}
			defer file.Close()
			file.WriteString(sourceCode)

			fmt.Println(color.GreenString("Saved submission to:"), fileName)

			// Move back to root directory.
			if err := os.Chdir(currentDir); err != nil {
				fmt.Println(color.RedString("unexpected error occurred:"), err)
				os.Exit(1)
			}
		}
	}
}

func shouldPull(sub codeforces.Submission, mode string, mp map[codeforces.Args]bool) bool {
	switch mode {
	case "l": // latest
		if _, ok := mp[sub.Arg]; ok == false {
			mp[sub.Arg] = true
			return true
		}
	case "la": // latest ac
		if sub.Verdict == "Accepted" {
			if _, ok := mp[sub.Arg]; ok == false {
				mp[sub.Arg] = true
				return true
			}
		}
	case "a": // all
		mp[sub.Arg] = true
		return true
	case "aa": // all ac
		if sub.Verdict == "Accepted" {
			mp[sub.Arg] = true
			return true
		}
	}

	return false
}
