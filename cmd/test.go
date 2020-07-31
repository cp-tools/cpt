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
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(testCmd)
	// define flags here
	testCmd.Flags().StringP("checker", "c", "lcmp", "Select (testlib) checker to use")
	testCmd.RegisterFlagCompletionFunc("checker", func(cmd *cobra.Command,
		_ []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		allCheckers := getCheckers(toComplete)
		for i, checker := range allCheckers {
			hCmd := exec.Command(checker, "--help")
			var data strings.Builder
			hCmd.Stderr = &data
			hCmd.Run()

			rgx := regexp.MustCompile(`Checker name: "([ -~]*)"`)
			if tmp := rgx.FindStringSubmatch(data.String()); len(tmp) > 1 {
				allCheckers[i] = fmt.Sprintf("%v\t%v", filepath.Base(allCheckers[i]), tmp[1])
				continue
			}
			baseName := filepath.Base(allCheckers[i])
			allCheckers[i] = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		}
		return allCheckers, cobra.ShellCompDirectiveNoFileComp
	},
	)

	testCmd.Flags().StringP("file", "f", "", "Solution file to run tests")
	testCmd.Flags().DurationP("time-limit", "t", time.Second, "Time limit for each test case")
	testCmd.Flags().BoolP("custom-invocation", "C", false, "Run solution, with input from stdin")
	testCmd.Flags().StringSlice("input", nil, "Test case inputs (corresponding to --output slice)")
	testCmd.Flags().StringSlice("output", nil, "Test case outputs (corresponding to --input slice)")

	testCmd.RunE = func(cmd *cobra.Command, _ []string) error {
		lflags := testCmd.Flags()

		if lflags.Changed("checker") && lflags.Changed("custom-invocation") {
			// both '--checker' and '--custom-invocation' specified
			return fmt.Errorf("Invalid flags - can't use both 'checker' and 'custom-invocation'")
		}

		if lflags.Changed("custom-invocation") && (lflags.Changed("input") || lflags.Changed("output")) {
			// both ('--input' or '--output') and '--custom-invocation' specified
			return fmt.Errorf("Invalid flags - can't use both 'input'/'output' and 'custom-invocation'")
		}

		inFiles, _ := lflags.GetStringSlice("input")
		outFiles, _ := lflags.GetStringSlice("output")
		if len(inFiles) != len(outFiles) {
			// lengths of slices '--input' and '--output' don't match
			return fmt.Errorf("Invalid flags - len of 'input' [%d] != len of 'output' [%d]", len(inFiles), len(outFiles))
		}

		// move this to utils function, getSampleTestFiles()....
		for _, file := range append(inFiles, outFiles...) {
			// check if all provided files exist
			if _, err := os.Stat(file); err != nil {
				return fmt.Errorf("Invalid flags - file %v doesn't exist", file)
			}
		}

		checker, _ := lflags.GetString("checker")
		allCheckers := getCheckers(checker)
		if len(allCheckers) == 0 {
			// no checker of value '--checker' exists (in checkers)
			return fmt.Errorf("Invalid flags - checker %v not found", checker)
		}
		lflags.Lookup("checker").Value.Set(allCheckers[0])

		// finally, run the test command
		test(lflags)
		return nil
	}
}

// return all matching checkers from ...cpt/checkers/
func getCheckers(toComplete string) []string {
	// find all checkers in checkers directory
	checkerDir := filepath.Join(cfgDir, "checkers")
	var allCheckers []string
	files, err := ioutil.ReadDir(checkerDir)
	if err != nil {
		// unexpected, what do I do?
		return nil
	}

	for _, file := range files {
		if file.Mode()&0111 == 0 || !strings.HasPrefix(file.Name(), toComplete) {
			// continue if not executable file matching description
			continue
		}
		allCheckers = append(allCheckers, filepath.Join(checkerDir, file.Name()))
	}
	return allCheckers
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
		tmplMap := make(map[string]interface{})
		tmplMap["file"] = file
		tmplMap["fileBasename"] = filepath.Base(file)
		tmplMap["fileBasenameNoExt"] = filepath.Ext(filepath.Base(file))

		err = tmpl.Execute(&scp, tmplMap)

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

			// @todo Run tests in parallel
			// @body Will speed up evaluation somewhat.
			// @body Tho the order of tests might not hold.

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

						inFile, err := os.Open(inFiles[i])
						if err != nil {
							panic(err)
						}
						defer inFile.Close()

						outFile, err := os.Open(outFiles[i])
						if err != nil {
							panic(err)
						}
						defer outFile.Close()

						oufFile, err := os.Open(oufFile.Name())
						if err != nil {
							panic(err)
						}
						defer oufFile.Close()

						inBuf := make([]byte, 70)
						if n, _ := inFile.Read(inBuf); n == 70 {
							// append '...' to the end
							inBuf = append(inBuf, []byte("...")...)
						}

						outBuf := make([]byte, 70)
						if n, _ := outFile.Read(outBuf); n == 70 {
							// append '...' to the end
							outBuf = append(outBuf, []byte("...")...)
						}

						oufBuf := make([]byte, 70)
						if n, _ := oufFile.Read(oufBuf); n == 70 {
							// append '...' to the end
							oufBuf = append(oufBuf, []byte("...")...)
						}

						diff := util.Diff(string(oufBuf), string(outBuf))
						tmplMap["diff"] = diff
						tmplMap["inp"] = string(inBuf)
					} else {
						// feel better yet?!
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
		tmplMap := make(map[string]interface{})
		tmplMap["file"] = file
		tmplMap["fileBasename"] = filepath.Base(file)
		tmplMap["fileBasenameNoExt"] = filepath.Ext(filepath.Base(file))

		err = tmpl.Execute(&scp, tmplMap)

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
