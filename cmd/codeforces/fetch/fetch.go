package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt-lib/v2/codeforces"
	"github.com/cp-tools/cpt/cmd/codeforces/open"
	"github.com/cp-tools/cpt/cmd/generate"
	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

	"github.com/fatih/color"
)

// Fetch fetches
func Fetch(arg codeforces.Args, cnf *conf.Conf) {
	// Extract contest (countdown) details.
	fmt.Println(color.BlueString("Fetching contest details of:"), arg)
	countdownDur, err := arg.GetCountdown()
	if err != nil {
		fmt.Println(color.RedString("error occurred while fetching contest details:"), err)
		os.Exit(1)
	}

	// Start countdown timer if contest has not started.
	if countdownDur.Seconds() > 0 {
		utils.RunCountdown(countdownDur, color.BlueString("Contest begins in:"))
		// Open problems page and dashboard once countdown done.
		open.Open(arg, "p")
		open.Open(arg, "d")
	}

	// Fetch required problems from contest page.
	problems, err := arg.GetProblems()
	if err != nil {
		fmt.Println(color.RedString("error occurred while fetching problems:"), err)
		os.Exit(1)
	}

	// Template to use to specify folder path for each fetched problem.
	folderPath := filepath.Join(cnf.GetStrings("fetch.problemFolderPath")...)
	if _, err := utils.CleanTemplate(folderPath, nil); err != nil {
		fmt.Println(color.RedString("error occurred while parsing 'folderPath' template:"), err)
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(color.RedString("unexpected error occurred:"), err)
		os.Exit(1)
	}

	for _, problem := range problems {
		// Determine folder path to parse problem to.
		problemDir, _ := utils.CleanTemplate(folderPath, problem)
		problemDir = filepath.Clean(problemDir)
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
		// Create sample tests files.
		testInFiles, testOutFiles := createTests(problem)
		// Create problem configuration (meta.yaml) file.
		problemCnf := createConf(problem, testInFiles, testOutFiles).SetParent(cnf)

		// Generate template if auto generation set to true.
		if cnf.GetBool("generate.onFetch") == true {
			alias := cnf.GetString("generate.defaultTemplate")
			if alias != "" && cnf.Has("template."+alias) {
				generate.Generate(alias, problemCnf, problemCnf.GetAll())
			}
		}

		// Move back to root directory.
		if err := os.Chdir(currentDir); err != nil {
			fmt.Println(color.RedString("unexpected error occurred:"), err)
			os.Exit(1)
		}
	}
}
