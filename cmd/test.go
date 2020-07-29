package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	testCmd.Flags().StringP("checker", "c", "lcmp", "Select output checker to use")
	testCmd.RegisterFlagCompletionFunc("checker",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}
			checkers := getCheckers(toComplete)
			for i, checker := range checkers {
				hCmd := exec.Command(checker, "--help")
				var data strings.Builder
				hCmd.Stderr = &data
				hCmd.Run()

				rgx := regexp.MustCompile(`Checker name: "([ -~]*)"`)
				if tmp := rgx.FindStringSubmatch(data.String()); len(tmp) > 1 {
					checkers[i] = fmt.Sprintf("%v\t%v", filepath.Base(checkers[i]), tmp[1])
					continue
				}
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

		if lflags.Changed("checker") && lflags.Changed("input") {
			return fmt.Errorf("Invalid flags - can't specify input/output files along with 'custom-invocation'")
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
	if script := tmplt["script"].(string); true {
		cmds, err := shellquote.Split(script)
		if err != nil {
			panic(err)
		}

		if ci, _ := lflags.GetBool("custom-invocation"); ci == true {
			// run custom invocation judge
			cmd := exec.Command(cmds[0], cmds[1:]...)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

			fmt.Println("-------START-------")
			cmd.Run()
			fmt.Println("--------END--------")

			fmt.Println()
		} else {
			// run checker judge (non interactive)
			inFiles, _ := lflags.GetStringSlice("input")
			outFiles, _ := lflags.GetStringSlice("output")
			util.FindInpOutFiles(&inFiles, &outFiles)

			timeLimitDur, _ := lflags.GetDuration("time-limit")

			// set verdict template data to parse
			tmplStr := "Test: #{{.testIndex}} -- Verdict: {{.verdict}} -- Time: {{.dur}}\n" +
				"------------------------------------------\n" +
				"{{- if .stderr}}\nStderr: {{.stderr}}{{end}}\n" +
				"{{- if eq .verdict \"WA\"}}\nInput\n{{.inp}}\n{{.diff}}{{end}}\n" +
				"{{- if .checkerLog}}\nChecker log: {{.checkerLog}}{{end}}\n"
			tmpl, _ := template.New("verdict").Parse(tmplStr)

			// run test for each input/output sample file(s)
			for i := 0; i < len(inFiles); i++ {
				// holds tmpl data values
				tmplMap := make(map[string]interface{})
				tmplMap["testIndex"] = i + 1

				ctx, cancel := context.WithTimeout(context.Background(), timeLimitDur)
				defer cancel()

				cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
				var cmdStdout, cmdStderr bytes.Buffer
				cmd.Stdout, cmd.Stderr = &cmdStdout, &cmdStderr

				inFile, err := os.Open(inFiles[i])
				if err != nil {
					fmt.Println("Could not read file", inFiles[i])
					fmt.Println("Skipping test case...")
					continue
				}
				defer inFile.Close()
				cmd.Stdin = inFile

				// run the executable now
				cmdStart := time.Now()
				err = cmd.Run()
				cmdTime := time.Since(cmdStart)
				cmdTime = cmdTime.Truncate(time.Millisecond)
				// write data to verdict tmpl
				tmplMap["stderr"] = cmdStderr.String()
				tmplMap["dur"] = cmdTime.String()

				select {
				case <-ctx.Done():
					// TLE timeout took place
					tmplMap["verdict"] = "TLE"
				default:
					// not a TLE, continue
					if err != nil {
						tmplMap["verdict"] = "RTE"
						break
					}

					oufFile, err := ioutil.TempFile(os.TempDir(), "ouf")
					defer os.Remove(oufFile.Name())
					if err != nil {
						// or should I not panic?
						panic(err)
					}
					if _, err = oufFile.Write(cmdStdout.Bytes()); err != nil {
						panic(err)
					}

					// run checker (testlib (args) - <input> <ouf> <out>)
					checker, _ := lflags.GetString("checker")
					var checkerStderr bytes.Buffer
					checkerCmd := exec.Command(checker, inFiles[i], oufFile.Name(), outFiles[i])
					checkerCmd.Stderr = &checkerStderr

					// run checker here
					err = checkerCmd.Run()
					if _, ok := err.(*exec.ExitError); ok {
						// there was an exit-error here!
						tmplMap["verdict"] = "WA"

						inBuf, _ := ioutil.ReadFile(inFiles[i])
						tmplMap["inp"] = string(inBuf)
						// diff prints all lines (todo: add '--no-diff' flag)
						outFile, err := os.Open(outFiles[i])
						defer outFile.Close()
						if err != nil {
							panic(err)
						}
						outBuf, _ := ioutil.ReadFile(outFiles[i])
						oufBuf, _ := ioutil.ReadFile(oufFile.Name())
						tmplMap["diff"] = util.DiffString(string(oufBuf), string(outBuf))
					} else {
						tmplMap["verdict"] = "AC"
					}
					tmplMap["checkerLog"] = checkerStderr.String()

				}
				tmpl.Execute(os.Stdout, tmplMap)
			}
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
Checker log: SUCCESS
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
