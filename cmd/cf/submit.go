package cf

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// submitCmd refers to 'submit' command
var submitCmd = &cobra.Command{
	Use:   "submit [SPECIFIER]",
	Short: "Submit solution file to problem",
}

func init() {
	RootCmd.AddCommand(submitCmd)

	// set flags in command
	submitCmd.Flags().StringP("file", "f", "", "Select file to submit")

	// set run command
	submitCmd.Run = func(cmd *cobra.Command, args []string) {
		lflags := submitCmd.Flags()

		// checking if file is valid is not reqd here
		// since it's done below while reading the file

		spfr, _ := util.DetectSpfr(args)
		submit(spfr, lflags)
	}
}

func submit(spfr string, lflags *pflag.FlagSet) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(arg.Contest) == 0 {
		fmt.Println("Contest id not specified")
		os.Exit(1)
	} else if len(arg.Problem) == 0 {
		fmt.Println("Problem id not specified")
		os.Exit(1)
	}

	// find all code files in current directory (check if given 'file' is valid)
	file, _ := lflags.GetString("file")
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
		// current user is loaded correctly
		fmt.Println("Current user:", handle)
	}

	langName := viper.GetString("templates." + tmpltAlias + ".languages.codeforces")
	err = arg.SubmitSolution(codeforces.LanguageID[langName], string(sourceData))
	if err != nil {
		fmt.Println("Could not submit code")
		fmt.Printf("%+q\n", err)
		os.Exit(1)
	}
	fmt.Println("Submitted")

	// watch submission row, and print data
	writer := uilive.New()
	writer.Start()
	for isDone := false; isDone == false; {
		start := time.Now()

		submissions, err := arg.GetSubmissions("")
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
