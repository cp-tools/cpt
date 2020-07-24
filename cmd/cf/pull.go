package cf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [SPECIFIER]",
	Short: "Pulls and saves submission code to local file",
}

func init() {
	RootCmd.AddCommand(pullCmd)

	// set flags in command
	var usernameFlag string
	pullCmd.Flags().StringVarP(&usernameFlag, "username", "u", "", "Username to fetch submissions of")
	/*var currDirFlag bool
	pullCmd.Flags().BoolVar(&currDirFlag, "current-directory", false, "Save code file(s) to current directory")
	var subIDFlag string
	pullCmd.Flags().StringVarP(&subIDFlag, "submission-id", "i", "", "Submission id of submission")*/

	pullCmd.RunE = func(cmd *cobra.Command, args []string) error {
		/*if _, err := strconv.Atoi(subIDFlag); len(subIDFlag) > 0 && err != nil {
			return fmt.Errorf("submission-id requires an integer value")
		}*/

		spfr, workDir := util.DetectSpfr(args)
		pull(spfr, workDir, usernameFlag)
		return nil
	}
}

func pull(spfr, workDir, usernameFlag string) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Fetching submission(s) details... Please wait!")
	submissions, err := arg.GetSubmissions(usernameFlag)
	if err != nil {
		fmt.Println("Could not pull submission(s) details")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(submissions) == 0 {
		fmt.Println("No submissions found")
		os.Exit(0)
	}

	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pMap := make(map[codeforces.Args]bool)
	var failPull []codeforces.Args
	// currently, only parses latest, AC submission.
	for _, sub := range submissions {
		if sub.Verdict != "Accepted" {
			continue
		}

		if val := pMap[sub.Arg]; val == true {
			continue
		}

		pMap[sub.Arg] = true

		fmt.Println("Pulling submission:", sub.ID)
		sourceCode, err := sub.GetSourceCode()
		if err != nil {
			fmt.Println("Could not pull submission")
			fmt.Println(err)
			// or should I break immediately?
			failPull = append(failPull, sub.Arg)
			continue
		}

		probDir := filepath.Join(workDir, "codeforces", sub.Arg.Class)
		if sub.Arg.Class == codeforces.ClassGroup {
			probDir = filepath.Join(probDir, sub.Arg.Group)
		}
		probDir = filepath.Join(probDir, sub.Arg.Contest, sub.Arg.Problem)
		err = os.MkdirAll(probDir, os.ModePerm)
		if err != nil {
			fmt.Println("Could not create problem folder")
			fmt.Println(err)
			os.Exit(1)
		}

		os.Chdir(probDir)
		// generate code file
		fileBase := sub.Arg.Problem
		fileExt := codeforces.LanguageExtn[sub.Language]
		for fileName, c := fileBase+fileExt, 1; true; c++ {
			if _, err := os.Stat(fileName); os.IsNotExist(err) == false {
				fmt.Println("File", fileName, "already exists in directory")
			} else {
				err = ioutil.WriteFile(fileName, []byte(sourceCode), 0644)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println("Saved submission", sub.Arg.Contest, sub.Arg.Problem, "to:", fileName)
				break
			}
			fileName = fmt.Sprintf("%v_%d%v", fileBase, c, fileExt)
		}
		os.Chdir(currDir)
	}

	// what do I do with failPull[] ??
}
