package cf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var pullCmd = &cobra.Command{
	Use:   "pull [SPECIFIER]",
	Short: "Pulls and saves submission code to local file",
}

func init() {
	RootCmd.AddCommand(pullCmd)

	// set flags in command
	pullCmd.Flags().StringP("username", "u", "", "Username to fetch submissions of")

	pullCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lflags := pullCmd.Flags()

		// set current user username if not set
		if !lflags.Changed("username") {
			usr := cfViper.GetString("username")
			if usr == "" {
				return fmt.Errorf("Invalid flags - 'username' not specified")
			}
			lflags.Lookup("username").Value.Set(cfViper.GetString("username"))
		}

		spfr, workDir := util.DetectSpfr(args)
		pull(spfr, workDir, lflags)
		return nil
	}
}

func pull(spfr, workDir string, lflags *pflag.FlagSet) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	color.Blue("Fetching submission details")
	username, _ := lflags.GetString("username")
	submissions, err := arg.GetSubmissions(username)
	if err != nil {
		color.Red("Could not pull submission(s) details")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(submissions) == 0 {
		color.Yellow("No submissions found")
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

		fmt.Println(color.BlueString("Pulling submission:"), sub.ID)
		sourceCode, err := sub.GetSourceCode()
		if err != nil {
			color.Red("Could not pull submission")
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
			color.Red("Could not create problem folder")
			fmt.Println(err)
			os.Exit(1)
		}

		os.Chdir(probDir)
		// generate code file
		fileBase := sub.Arg.Problem
		fileExt := codeforces.LanguageExtn[sub.Language]
		for fileName, c := fileBase+fileExt, 1; true; c++ {
			if _, err := os.Stat(fileName); os.IsNotExist(err) == false {
				color.Yellow("File %v already exists in directory", fileName)
			} else {
				err = ioutil.WriteFile(fileName, []byte(sourceCode), 0644)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				color.Green("Saved submission in: %v", sub.Arg)
				break
			}
			fileName = fmt.Sprintf("%v_%d%v", fileBase, c, fileExt)
		}
		os.Chdir(currDir)
	}
	// what do I do with failPull[] ??
}
