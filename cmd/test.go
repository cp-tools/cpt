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
	"sync"
	"text/template"
	"time"

	"github.com/cp-tools/cpt/util"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
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

		inFiles := lflags.MustGetStringSlice("input")
		outFiles := lflags.MustGetStringSlice("output")
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

		checker := lflags.MustGetString("checker")
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
	file, err := util.FindCodeFiles(lflags.MustGetString("file"))
	if err != nil {
		color.Red("Could not select code file")
		fmt.Println(err)
		os.Exit(1)
	}
	// find template configuration to use
	tmpltAlias, err := util.FindTemplateToUse(file)
	if err != nil {
		color.Red("Could not select template configuration")
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

		fmt.Println(color.BlueString("Prescript:"), scp.String())
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

	// run script next
	if script := tmplt["script"].(string); true {
		cmds, err := shellquote.Split(script)
		if err != nil {
			panic(err)
		}

		switch {
		case lflags.MustGetBool("custom-invocation"):
			param := make(map[string]interface{})
			param["cmds"] = cmds
			customInvocation(param)

		default:
			// non interactive judge to run
			inFilesPath := lflags.MustGetStringSlice("input")
			outFilesPath := lflags.MustGetStringSlice("output")
			util.FindInpOutFiles(&inFilesPath, &outFilesPath)

			tlDur := lflags.MustGetDuration("time-limit")

			fvt := color.New(color.FgBlue).Add(color.Bold).SprintFunc()
			// set verdict template to use
			tmplStr := fmt.Sprintf("%v #{{.testIndex}}\t%v {{.verdict}}\t%v {{.duration}}\n", fvt("Test:"), fvt("Verdict:"), fvt("Time:")) +
				fmt.Sprintf("{{- if .log}}\n%v {{.log}}{{end}}\n", fvt("Fail log:")) +
				fmt.Sprintf("{{- if .stderr}}\n%v {{.stderr}}{{end}}\n", fvt("Stderr:")) +
				fmt.Sprintf("{{- if .checkerLog}}\n%v {{.checkerLog}}{{end}}\n", fvt("Checker log:")) +
				fmt.Sprintf("{{- if .diff}}\n%v\n{{.input}}\n{{.diff}}{{end}}", util.HeaderCol("Input"))
			tmpl, _ := template.New("verdictTmpl").Option("missingkey=zero").Parse(tmplStr)

			var wg sync.WaitGroup
			for i := 0; i < len(inFilesPath); i++ {
				wg.Add(1)

				param := map[string]interface{}{
					"cmds":        cmds,
					"inFilePath":  inFilesPath[i],
					"outFilePath": outFilesPath[i],
					"time-limit":  tlDur,
					"checker":     lflags.MustGetString("checker"),
				}

				go func(i int) {
					defer wg.Done()
					tmplMap := checkerJudge(param)
					tmplMap["testIndex"] = i + 1
					var buf strings.Builder
					tmpl.Execute(&buf, tmplMap)
					str := strings.TrimSpace(buf.String())
					fmt.Println(str + "\n")
				}(i)
			}
			wg.Wait()
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

		fmt.Println(color.BlueString("Postscript:"), scp.String())
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

// checkerJudge is the checker based judge.
// Data to be passed are:
// cmds ([]string) => script commands split
// inFilePath (string) => input file path
// outFilePath (string) => output file path
// time-limit (time.Duration) => run timelimit
// checker (string) => path to checker file to use
func checkerJudge(param map[string]interface{}) (tmplMap map[string]interface{}) {
	tmplMap = make(map[string]interface{})

	// handle panic, recover
	defer func() {
		if r := recover(); r != nil {
			tmplMap["verdict"] = color.RedString("FAIL")
			tmplMap["log"] = r
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), param["time-limit"].(time.Duration))
	defer cancel()

	cmds := param["cmds"].([]string)
	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	var cmdStdout, cmdStderr bytes.Buffer
	cmd.Stdout = &cmdStdout
	cmd.Stderr = &cmdStderr

	inFile, err := os.Open(param["inFilePath"].(string))
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	cmd.Stdin = inFile

	// run the executable
	start := time.Now()
	err = cmd.Run()
	since := time.Since(start).Truncate(time.Millisecond)

	tmplMap["stderr"] = cmdStderr.String()
	tmplMap["duration"] = since.String()

	if since >= param["time-limit"].(time.Duration) {
		// Time Limit Exceeded occured.
		tmplMap["verdict"] = color.YellowString("TLE")
		return tmplMap
	}
	// not a TLE; continue
	if err != nil {
		// runtime error occured.
		tmplMap["verdict"] = color.RedString("RTE")
		return tmplMap
	}

	// write submission output to temp file
	oufFile, err := ioutil.TempFile(os.TempDir(), "ouf")
	if err != nil {
		panic(err)
	}
	defer os.Remove(oufFile.Name())
	if _, err = oufFile.Write(cmdStdout.Bytes()); err != nil {
		panic(err)
	}

	// run checker (<checker> <inp> <ouf> <out>)
	var checkerStderr bytes.Buffer
	checkerCmd := exec.Command(param["checker"].(string), param["inFilePath"].(string),
		oufFile.Name(), param["outFilePath"].(string))
	checkerCmd.Stderr = &checkerStderr

	err = checkerCmd.Run()

	tmplMap["checkerLog"] = string(checkerStderr.Bytes())
	if _, ok := err.(*exec.ExitError); ok {
		// there was a non-zero exit code. WA output
		tmplMap["verdict"] = color.RedString("WA")

		inFileData, err := ioutil.ReadFile(param["inFilePath"].(string))
		if err != nil {
			panic(err)
		}
		if len(inFileData) > 70 {
			inFileData = append(inFileData[:70], []byte("...")...)
		}

		outFileData, err := ioutil.ReadFile(param["outFilePath"].(string))
		if err != nil {
			panic(err)
		}
		if len(outFileData) > 70 {
			outFileData = append(outFileData[:70], []byte("...")...)
		}

		oufFileData, err := ioutil.ReadFile(oufFile.Name())
		if err != nil {
			panic(err)
		}
		if len(oufFileData) > 70 {
			oufFileData = append(oufFileData[:70], []byte("...")...)
		}

		diff := uitable.New()
		diff.Separator = " | "
		diff.AddRow(util.HeaderCol("OUTPUT"), util.HeaderCol("ANSWER"))
		diff.AddRow(string(oufFileData), string(outFileData))

		tmplMap["diff"] = diff.String()
		tmplMap["input"] = string(inFileData)
		return tmplMap
	} else if err != nil {
		panic(err)
	}

	tmplMap["verdict"] = color.GreenString("AC")
	return tmplMap
}

// customInvocation is for well, custom invocation run.
// Data to be passed are:
// cmds ([]string) => script commands split
func customInvocation(param map[string]interface{}) {
	cmds := param["cmds"].([]string)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	color.Green("---- * ---- * ---- * ----")
	cmd.Run()
	color.Green("---- * ---- * ---- * ----")
	fmt.Println()
}
