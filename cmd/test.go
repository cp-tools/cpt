package cmd

import (
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
	testCmd.Flags().StringP("checker", "c", "", "Select output checker to use")
	testCmd.RegisterFlagCompletionFunc("checker",
		func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getCheckers(toComplete), cobra.ShellCompDirectiveDefault
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
		flags := testCmd.Flags()

		if flags.Changed("checker") && flags.Changed("custom-invocation") {
			return fmt.Errorf("Invalid flags - can't use both 'checker' and 'custom-invocation'")
		}

		inpf, _ := flags.GetStringSlice("input")
		outf, _ := flags.GetStringSlice("output")
		if len(inpf) != len(outf) {
			// check if lengths of test case slice match
			return fmt.Errorf("Invalid flag values - len of 'input' [%d], doesn't match len of 'output' [%d]", len(inpf), len(outf))
		}
		for _, fp := range append(inpf, outf...) {
			// check if all provided files exist
			if _, err := os.Stat(fp); err != nil {
				return fmt.Errorf("File %v doesn't exist", fp)
			}
		}

		test(testCmd.Flags())
		return nil
	}
}

// return all matching checkers from ...cpt/checkers/
func getCheckers(toComplete string) []string {
	var checkers []string
	// read checkers from .../cpt/checkers/
	files, err := ioutil.ReadDir(filepath.Join(cfgDir, "checkers"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileBase := filepath.Base(file.Name())
		if file.IsDir() || file.Mode()&0111 == 0 || !strings.HasPrefix(fileBase, toComplete) {
			// continue if not executable file matching desc
			continue
		}

		checkers = append(checkers, fileBase)
	}
	cobra.CompDebugln(strings.Join(checkers, " "), false)
	return checkers
}

func test(flags *pflag.FlagSet) {
	// find code file to test
	file, _ := flags.GetString("file")
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

	if ci, _ := flags.GetBool("custom-invocation"); ci == true {
		cmds, err := shellquote.Split(tmplt["script"].(string))
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("-------START-------")
		cmd.Run()
		fmt.Println("--------END--------")
	} else {
		// find all input/output files
		inpf, _ := flags.GetStringSlice("input")
		outf, _ := flags.GetStringSlice("output")
		inpf, outf = util.FindInpOutFiles(inpf, outf)

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
