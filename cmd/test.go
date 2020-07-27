package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/cp-tools/cpt/util"

	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run solution file against test cases",
}

func init() {
	rootCmd.AddCommand(testCmd)
	// define flags here
	testCmd.Flags().StringP("checker", "c", "basic", "Select output checker to use")
	testCmd.RegisterFlagCompletionFunc("checker",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}
			checkers := getCheckers(toComplete)
			for i := range checkers {
				checkers[i] = filepath.Base(checkers[i])
			}
			return checkers, cobra.ShellCompDirectiveNoFileComp
		},
	)

	testCmd.Flags().StringP("file", "f", "", "Solution file to run tests")
	testCmd.Flags().DurationP("time-limit", "t", time.Second, "Time limit for each test case")
	testCmd.Flags().BoolP("custom-invocation", "C", false, "Run solution, with input from stdin")
	testCmd.Flags().StringSlice("input", nil, "Test case inputs (corresponding to --output slice)")
	testCmd.Flags().StringSlice("output", nil, "Test case outputs (corresponding to --input slice)")

	// the execution part. idk why I'm even saying this :-\
	// maybe cause it's 4 at night? idk XD
	testCmd.RunE = func(cmd *cobra.Command, _ []string) error {
		lflags := testCmd.Flags()

		if lflags.Changed("checker") && lflags.Changed("custom-invocation") {
			return fmt.Errorf("Invalid flags - can't use both 'checker' and 'custom-invocation'")
		}

		inpf, _ := lflags.GetStringSlice("input")
		outf, _ := lflags.GetStringSlice("output")
		if len(inpf) != len(outf) {
			// check if lengths of test case slice match
			return fmt.Errorf("Invalid flag values - len of 'input' [%d], doesn't match len of 'output' [%d]", len(inpf), len(outf))
		}
		for _, fp := range append(inpf, outf...) {
			// check if all provided files exist
			if _, err := os.Stat(fp); err != nil {
				return fmt.Errorf("Invalid test files - file %v doesn't exist", fp)
			}
		}
		checker, _ := lflags.GetString("checker")
		if len(getCheckers(checker)) == 0 {
			return fmt.Errorf("Invalid checker - checker %v not found in %v",
				checker, filepath.Join(cfgDir, "cpt", "checkers", "."))
		}
		lflags.Lookup("checker").Value.Set(getCheckers(checker)[0])

		test(lflags)
		return nil
	}
}

// return all matching checkers from ...cpt/checkers/
func getCheckers(toComplete string) []string {
	var checkers []string
	// read checkers from .../cpt/checkers/
	ckrDir := filepath.Join(cfgDir, "checkers")
	files, err := ioutil.ReadDir(ckrDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() || file.Mode()&0111 == 0 || !strings.HasPrefix(file.Name(), toComplete) {
			// continue if not executable file matching desc
			continue
		}

		checkers = append(checkers, filepath.Join(ckrDir, file.Name()))
	}
	return checkers
}

func test(lflags *pflag.FlagSet) {
	// find code file to test
	file, _ := lflags.GetString("file")
	file, err := util.FindCodeFiles(file)
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

	tmplt := viper.GetStringMap("templates." + tmpltAlias)

	// run pre script first
	if script := tmplt["prescript"].(string); len(script) > 0 {
		tmpl, err := template.New("prescript").Parse(script)
		if err != nil {
			panic(err)
		}

		var scp strings.Builder
		err = tmpl.Execute(&scp, map[string]string{
			"file": file,
		})

		fmt.Println("Prescript:", scp.String())
		cmds, err := shellquote.Split(scp.String())
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println()
	}

	// run script
	sCmds, err := shellquote.Split(tmplt["script"].(string))

	if ci, _ := lflags.GetBool("custom-invocation"); ci == true {
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(sCmds[0], sCmds[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("-------START-------")
		cmd.Run()
		fmt.Println("--------END--------")
	} else {
		// find all input/output files
		inpf, _ := lflags.GetStringSlice("input")
		outf, _ := lflags.GetStringSlice("output")
		inpf, outf = util.FindInpOutFiles(inpf, outf)

		dur, _ := lflags.GetDuration("time-limit")

		verdictTmpl, _ := template.New("verdict").Parse(
			`Test: #{{.testIndex}} -- Verdict: {{.verdict}} -- Time: {{.dur}}
{{- if .stderr}}
Stderr: {{.stderr}}
{{- end}}
{{- if eq .verdict "WA"}}
Input
{{.inp}}
{{.diff}}
{{- end}}
{{- if .checkerLog }}
Checker log: {{.checkerLog}}
{{- end}}{{- /* remove extra space at the end */}}
--------------------------------------------`)

		// run tests for each sample case
		for idx := 0; idx < len(inpf); idx++ {
			ctx, cancel := context.WithTimeout(context.Background(), dur)
			defer cancel()
			cmd := exec.CommandContext(ctx, sCmds[0], sCmds[1:]...)

			// holds verdictTmpl data values

			var cmdStdoutT, cmdStderrT strings.Builder
			cmd.Stdout, cmd.Stderr = &cmdStdoutT, &cmdStderrT
			inpData, err := ioutil.ReadFile(inpf[idx])
			if err != nil {
				panic(err)
			}

			cmd.Stdin = bytes.NewReader(inpData)

			start := time.Now()
			err = cmd.Run()
			since := time.Since(start)

			tplData := map[string]interface{}{
				"testIndex": idx + 1,
				"stderr":    cmdStderrT.String(),
				"dur":       since.Truncate(time.Millisecond).String(),
			}

			if since >= dur {
				tplData["verdict"] = "TLE"
			} else if err != nil {
				tplData["verdict"] = "RTE"
			} else {
				// create ouf temporary file
				oufFile, err := ioutil.TempFile(os.TempDir(), "ouf")
				if err != nil {
					panic(err)
				}
				oufFile.WriteString(cmdStdoutT.String())

				// run checker
				var checkerStdoutT strings.Builder
				checker, _ := lflags.GetString("checker")
				ckrCmd := exec.Command(checker, "--ans", outf[idx], "--ouf", oufFile.Name())
				ckrCmd.Stdout = &checkerStdoutT
				ckrCmd.Stderr = os.Stderr

				if err := ckrCmd.Run(); err != nil {
					if exitError, ok := err.(*exec.ExitError); ok {
						if exitError.ExitCode() != 1 {
							panic(err)
						}
					}
					// read output file
					outData, err := ioutil.ReadFile(outf[idx])
					if err != nil {
						panic(err)
					}

					// returns 1 means its a WA
					tplData["verdict"] = "WA"
					tplData["inp"] = string(inpData)
					tplData["diff"] = util.DiffString(cmdStdoutT.String(), string(outData))
					// add diff here
				} else {
					// it's an AC! Feel better already?
					tplData["verdict"] = "AC"
				}
				tplData["checkerLog"] = strings.TrimSpace(checkerStdoutT.String())
				defer os.Remove(oufFile.Name())
			}

			var verdict strings.Builder
			verdictTmpl.Execute(&verdict, tplData)
			fmt.Println(verdict.String())
		}
	}

	// run post script here (last)
	if script := tmplt["postscript"].(string); len(script) > 0 {
		tmpl, err := template.New("postscript").Parse(script)
		if err != nil {
			panic(err)
		}

		var scp strings.Builder
		err = tmpl.Execute(&scp, map[string]string{
			"file": file,
		})

		fmt.Println()
		fmt.Println("Postscript:", scp.String())
		cmds, err := shellquote.Split(scp.String())
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}

/*
Test: #1 -- Verdict: OK -- Time: 175ms
Checker: SUCCESS
--------------------------------------
Test: #2 -- Verdict: WA -- Time: 45ms
Checker log: Expected 'YES', found 'NO' (term 2)
Output		| Expected output
Yes			| Yes
No			| Yes
No			| No
Yes			| No
--------------------------------------
Test: #3 -- Verdict: TLE -- Time: 0ms
Stderr: This app is bugged
Don't ask why, but it's hacky
--------------------------------------
*/
