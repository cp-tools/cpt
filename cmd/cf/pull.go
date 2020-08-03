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
	Long: `Pulls submissions of specified user to local copy.
Pulls latest AC submission of user to directory similar to 'cpt cf fetch'.
By default, pulls the submission of current active user. Use flag --username
to pull submissions of particular user.

If SPECIFIER is not given, pulls all submissions matching above criteria.
Please DON'T MISUSE THIS TO DDOS CODEFORCES, as it is resource intensive.

Refer 'cpt cf -h' for details on argument [SPECIFIER].

Usage examples:
cpt cf pull 4 a
                            Pulls submissions of current user in problem 4a
cpt cf pull -u cp-tools
                            Pulls all submissions of user 'cp-tools'
`,
}

func init() {
	RootCmd.AddCommand(pullCmd)

	// set flags in command
	pullCmd.Flags().StringP("username", "u", "", "Username to fetch submissions of")

	pullCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lflags := pullCmd.Flags()

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
	submissions, err := arg.GetSubmissions(lflags.MustGetString("username"))
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
