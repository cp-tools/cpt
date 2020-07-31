package cf

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cp-tools/cpt-lib/codeforces"
	"github.com/cp-tools/cpt/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var listCmd = &cobra.Command{
	Use:       "list [MODE] [SPECIFIER]",
	Short:     "Prints tabulated results of various real-time data",
	ValidArgs: []string{"submissions", "dashboard", "contests"},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// set flags here
	listCmd.Flags().UintP("number", "n", 5, "Number of rows to output [contests submissions]")
	listCmd.Flags().BoolP("register", "r", false, "Enable registration menu [contests]")
	listCmd.Flags().StringP("username", "u", "", "Username to fetch data of [submissions]")

	// set listCmd Args validations
	listCmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Invalid args - requires mode argument")
		}
		if util.SliceContains(args[0], listCmd.ValidArgs) {
			return nil
		}
		return fmt.Errorf("Invalid args - mode '%v' not valid", args[0])
	}
	// set listCmd Run command
	listCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lflags := listCmd.Flags()
		// various flag combination validators
		switch args[0] {
		case "contests":
			if lflags.Changed("username") {
				// can't use username with contests (arg)
				return fmt.Errorf("Invalid flags - 'username' not applicable for mode 'contests'")
			}
		case "dashboard":
			if lflags.Changed("username") || lflags.Changed("number") || lflags.Changed("number") {
				// can't use username any flag with dashboard
				return fmt.Errorf("Invalid flags - mode 'dashboard' takes no flags")
			}
		case "submissions":
			if lflags.Changed("register") {
				// can't use register with submissions (arg)
				return fmt.Errorf("Invalid flags - 'register' not applicable for mode 'dashboard'")
			}
			// set current user username if not set
			if !lflags.Changed("username") {
				username := cfViper.GetString("username")
				if username == "" {
					return fmt.Errorf("Invalid flags - 'username' not specified")
				}
				lflags.Lookup("username").Value.Set(cfViper.GetString("username"))
			}
		}

		spfr, _ := util.DetectSpfr(args[1:])
		list(spfr, args[0], lflags)
		return nil
	}
}

func list(spfr, mode string, lflags *pflag.FlagSet) {
	arg, err := codeforces.Parse(spfr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch mode {
	case "submissions":
		// watch subissions row
		writer := uilive.New()
		writer.Start()
		for isJudging := false; ; isJudging = false {
			start := time.Now()

			username, _ := lflags.GetString("username")
			submissions, err := arg.GetSubmissions(username)
			if err != nil {
				color.Red("Could not fetch submissions")
				fmt.Println(err)
				os.Exit(1)
			}

			if len(submissions) == 0 {
				color.Yellow("No submissions found")
				os.Exit(0)
			}

			t := uitable.New()
			t.Separator = " | "
			t.MaxColWidth = 22

			t.AddRow(util.HeaderCol("#"), util.HeaderCol("When"), util.HeaderCol("Problem"), util.HeaderCol("Lang"),
				util.HeaderCol("Verdict"), util.HeaderCol("Time"), util.HeaderCol("Memory"))

			for i, sub := range submissions {
				number, _ := lflags.GetUint("number")
				if uint(i) >= number {
					break
				}
				if sub.IsJudging == true {
					isJudging = true
				}

				// compress verdict string
				verdict := sub.Verdict
				verdict = strings.ReplaceAll(verdict, "Time limit exceeded", "TLE")
				verdict = strings.ReplaceAll(verdict, "Compilation error", "CE")
				verdict = strings.ReplaceAll(verdict, "Runtime error", "RTE")
				verdict = strings.ReplaceAll(verdict, "Wrong answer", "WA")
				verdict = strings.ReplaceAll(verdict, "Accepted", "AC")
				verdict = strings.ReplaceAll(verdict, "Partial result", "PR")

				if strings.HasPrefix(verdict, "CE") || strings.HasPrefix(verdict, "TLE") {
					verdict = color.HiYellowString(verdict)
				} else if strings.HasPrefix(verdict, "RTE") || strings.HasPrefix(verdict, "WA") {
					verdict = color.HiRedString(verdict)
				} else if verdict == "AC" || strings.HasPrefix(verdict, "PR") {
					verdict = color.HiGreenString(verdict)
				}

				t.AddRow(sub.ID, sub.When.Local().Format("Jan/02/2006 15:04"), sub.Problem,
					sub.Language, verdict, sub.Time, sub.Memory)
			}
			fmt.Fprintln(writer, t.String())
			if isJudging == false {
				break
			}

			time.Sleep(time.Second*2 - time.Since(start))
		}
		writer.Stop()

	case "dashboard":
		dhbd, err := arg.GetDashboard()
		if err != nil {
			color.Red("Could not fetch dashboard")
			fmt.Println(err)
			os.Exit(1)
		}

		// list contest name
		fmt.Println(color.BlueString("Contest name:"), dhbd.Name)
		fmt.Println()

		// list countdown to contest end
		if dhbd.Countdown != 0 {
			fmt.Println(color.BlueString("Contest ends in:"), dhbd.Countdown.String())
			fmt.Println()
		}

		t := uitable.New()
		t.Separator = " | "
		t.MaxColWidth = 25

		t.AddRow(util.HeaderCol("Name"), util.HeaderCol("Status"), util.HeaderCol("Solved"))
		for _, prob := range dhbd.Problem {
			var solveStatus string
			switch prob.SolveStatus {
			case codeforces.SolveAccepted:
				solveStatus = color.HiGreenString("AC")
			case codeforces.SolveNotAttempted:
				solveStatus = "NA"
			case codeforces.SolveRejected:
				solveStatus = color.HiRedString("WA")
			}

			t.AddRow(prob.Name, solveStatus, prob.SolveCount)
		}
		fmt.Println(t.String())

	case "contests":
		// default to contests menu
		if len(arg.Class) == 0 {
			register, _ := lflags.GetBool("register")
			if register == true {
				// it means contests
				arg.Class = codeforces.ClassContest
			} else {
				util.SurveyErr(survey.AskOne(&survey.Select{
					Message: "Select contest class to list:",
					Options: []string{codeforces.ClassContest, codeforces.ClassGym},
					Default: codeforces.ClassContest,
				}, &arg.Class))
			}
		}

		var omitFinishedContests bool
		// omit finished contests if contest not set
		if len(arg.Contest) == 0 {
			omitFinishedContests = true
		} else {
			omitFinishedContests = false
		}

		contests, err := arg.GetContests(omitFinishedContests)
		if err != nil {
			color.Red("Could not fetch contests")
			fmt.Println(err)
			os.Exit(1)
		}

		t := uitable.New()
		t.Separator = " | "
		t.MaxColWidth = 30

		t.AddRow(util.HeaderCol("#"), util.HeaderCol("Name"), util.HeaderCol("Writers"), util.HeaderCol("Start"),
			util.HeaderCol("Length"), util.HeaderCol("Registration"), util.HeaderCol("Count"))
		for c, cont := range contests {
			number, _ := lflags.GetUint("number")
			if uint(c) >= number {
				break
			}

			var regStatus string
			switch cont.RegStatus {
			case codeforces.RegistrationOpen:
				regStatus = color.HiGreenString("OPEN")
			case codeforces.RegistrationClosed:
				regStatus = color.HiRedString("CLOSED")
			case codeforces.RegistrationDone:
				regStatus = color.HiGreenString("REGISTERED")
			case codeforces.RegistrationNotExists:
				regStatus = "NO REGISTRATION"
			}

			t.AddRow(cont.Arg.Contest, cont.Name, strings.Join(cont.Writers, " "),
				cont.StartTime.Local().Format("Jan/02/2006 15:04"),
				cont.Duration.String(), regStatus, cont.RegCount)
		}
		fmt.Println(t.String())

		// give user chance to register
		register, _ := lflags.GetBool("register")
		if register == true && arg.Class == codeforces.ClassContest {
			fmt.Println()

			var regOpenContestsName []string
			var regOpenContests []codeforces.Contest
			for c, cont := range contests {
				number, _ := lflags.GetUint("number")
				if uint(c) >= number {
					break
				}
				if cont.RegStatus == codeforces.RegistrationOpen {
					regOpenContests = append(regOpenContests, cont)
					regOpenContestsName = append(regOpenContestsName, cont.Name)
				}
			}

			var idxChoice int
			util.SurveyErr(survey.AskOne(&survey.Select{
				Message: "Select contest to register in:",
				Options: regOpenContestsName,
			}, &idxChoice))

			regInfo, err := regOpenContests[idxChoice].Arg.RegisterForContest()
			if err != nil {
				color.Red("Could not fetch registration page")
				fmt.Println(err)
				os.Exit(1)
			}

			var cfm bool
			util.SurveyErr(survey.AskOne(&survey.Confirm{
				Message: "Agree to terms and conditions (Enter '?' to view)?",
				Help:    regInfo.Terms,
				Default: false,
			}, &cfm))

			if cfm == false {
				color.Yellow("Registration aborted")
				os.Exit(0)
			}

			fmt.Println(color.BlueString("Registering in contest:"), regOpenContests[idxChoice].Arg.Contest)
			err = regInfo.Register()
			if err != nil {
				color.Red("Could not register user in contest")
				fmt.Println(err)
				os.Exit(1)
			}

			color.Green("Registered successfully!")
		}
	}
}
