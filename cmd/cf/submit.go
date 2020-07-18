package cf

import (
	"fmt"
	"os"
	"time"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

// submitCmd refers to 'submit' command
var submitCmd = &cobra.Command{
	Use: "submit [SPECIFIER]",
	Run: func(cmd *cobra.Command, args []string) {
		submit(util.DetectSpfr(args))
	},
}

func init() {
	RootCmd.AddCommand(submitCmd)

	submitCmd.Flags().StringP("file", "f", "", "Select file to submit")
	lFlags = submitCmd.Flags()
}

func submit(spfr, workDir string) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	/*
		if len(arg.Contest) == 0 {
			fmt.Println("Contest id not specified")
			os.Exit(1)
		} else if len(arg.Problem) == 0 {
			fmt.Println("Problem id not specified")
			os.Exit(1)
		}

		// find all code files in current directory, if file not specified
		file, _ := lFlags.GetString("file")
		file, err = util.FindCodeFiles(file)
		if err != nil {
			fmt.Println("Could not select code file")
			fmt.Println(err)
			os.Exit(1)
		}

		// find template configuration to use
		tmpltAlias, err := util.FindTemplateToUse(file)
		if err != nil {
			fmt.Println("Could not select template configuration")
			fmt.Println(err)
			os.Exit(1)
		}

		// read selected source file (and submit)
		sourceData, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Could not read code file")
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Checking login status")
		if usr, passwd := cfViper.GetString("username"), cfViper.GetString("password"); len(usr) == 0 || len(passwd) == 0 {
			fmt.Println("Could not find any saved login credentials")
			os.Exit(1)
		} else {
			passwd, err := util.Decrypt(usr, passwd)
			if err != nil {
				fmt.Println("Could not decrypt password")
				fmt.Println(err)
				os.Exit(1)
			}

			handle, err := codeforces.Login(usr, passwd)
			if err != nil {
				fmt.Println("Could not check login status")
				fmt.Println(err)
				os.Exit(1)
			}
			// current user is loaded correctly here
			fmt.Println("Current user:", handle)
			viper.WriteConfig()
		}

		tmplt := viper.GetStringMap("templates")[tmpltAlias].(map[string]interface{})
		langName := tmplt["languages"].(map[string]interface{})["codeforces"].(string)
		err = arg.SubmitSolution(codeforces.LanguageID[langName], string(sourceData))
		if err != nil {
			fmt.Println("Could not submit code")
			fmt.Printf("%+q\n", err)
			os.Exit(1)
		}
	*/

	// watch submission row, and print data
	writer := uilive.New()
	writer.Start()
	for isDone := false; isDone == false; {
		start := time.Now()

		submissions, err := arg.GetSubmissions(cfViper.GetString("username"), true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(submissions) == 0 {
			fmt.Println("Expected atleast 1 submission for problem, found 0")
			fmt.Println("Quiting....")
			os.Exit(1)
		}

		t := uitable.New()
		t.Separator = " "

		t.AddRow("Verdict:", submissions[0].Verdict)
		if submissions[0].IsJudging == false {
			t.AddRow("Memory:", submissions[0].Memory)
			t.AddRow("Time:", submissions[0].Time)
			isDone = true
		}

		fmt.Fprintln(writer, t.String())
		time.Sleep(time.Second - time.Since(start))
	}
	writer.Stop()
}
