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
	"github.com/spf13/viper"
)

var (
	// fetchCmd represents the 'cf fetch' command
	fetchCmd = &cobra.Command{
		Use:   "fetch [SPECIFIER]",
		Short: "Fetch and save sample tests from website to local folder",
		Run: func(cmd *cobra.Command, args []string) {
			fetch(util.DetectSpfr(args))
		},
	}

	// GenFunc to run 'cpt gen'
	// I hate cyclic dependencies
	GenFunc func(*pflag.FlagSet)
)

func init() {
	RootCmd.AddCommand(fetchCmd)
}

func fetch(spfr, workDir string) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(arg.Contest) == 0 {
		color.Red("Contest id not specified")
		os.Exit(1)
	}

	fmt.Println(color.BlueString("Fetching details of:"), arg)
	dur, err := arg.GetCountdown()
	if err != nil {
		color.Red("Could not extract countdown")
		fmt.Println(err)
		os.Exit(1)
	}

	if dur.Seconds() > 0 {
		util.StartCountdown(dur, color.BlueString("Contest starts in:"))
		// open problems page once parsing is done
		open(spfr)
	}

	// fetch all problems from contest page
	problems, err := arg.GetProblems()
	if err != nil {
		color.Red("Could not fetch sample tests")
		fmt.Println(err)
		os.Exit(1)
	}

	for _, prob := range problems {
		// set problem folder directory path
		probDir := filepath.Join(workDir, "codeforces", prob.Arg.Class)
		if prob.Arg.Class == codeforces.ClassGroup {
			probDir = filepath.Join(probDir, prob.Arg.Group)
		}
		probDir = filepath.Join(probDir, prob.Arg.Contest, prob.Arg.Problem)

		err := os.MkdirAll(probDir, os.ModePerm)
		if err != nil {
			color.Red("Could not create problem folder")
			fmt.Println(err)
			os.Exit(1)
		}

		// save test cases
		for c, sampleTest := range prob.SampleTests {
			inFilePath := filepath.Join(probDir, fmt.Sprintf("%d.in", c))
			outFilePath := filepath.Join(probDir, fmt.Sprintf("%d.out", c))

			err := ioutil.WriteFile(inFilePath, []byte(sampleTest.Input), os.ModePerm)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = ioutil.WriteFile(outFilePath, []byte(sampleTest.Output), os.ModePerm)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		fmt.Println(color.BlueString("Fetched %v sample tests:", len(prob.SampleTests)), prob.Name)

		// generate template if specified
		genTmplt := viper.GetString("default_template")
		if viper.GetBool("gen_on_fetch") && genTmplt != "none" {
			// set flags to run 'gen' command
			var genFlags pflag.FlagSet
			// do I have any other option?!
			genFlags.String("template", genTmplt, "")

			currDir, _ := os.Getwd()
			os.Chdir(probDir)
			GenFunc(&genFlags)
			os.Chdir(currDir)
		}
	}
}
