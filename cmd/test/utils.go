package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
)

func extractTestsFiles(cnf *conf.Conf) (inputFiles, expectedFiles []string) {
	inputFiles = cnf.GetStrings("problem.test.input")
	expectedFiles = cnf.GetStrings("problem.test.output")

	if len(inputFiles) != len(expectedFiles) {
		// Mismatch in test cases count.
		fmt.Println(color.RedString("error selecting test files:"), fmt.Sprintf("number of 'inputFiles' [%d] not equals number of 'expectedFiles' [%d]", len(inputFiles), len(expectedFiles)))
		os.Exit(1)
	}
	return
}

// SelectSubmissionFile returns template of submission file to use, based
// on configured templates and passed 'submissionFilePath' value.
// 'submissionFilePath' must point to a valid file, if specified.
func SelectSubmissionFile(submissionFilePath *string, cnf *conf.Conf) string {
	// Find all generated submission files in local configurations.
	generatedFilesMap := make(map[string]string)
	for _, templateAlias := range cnf.GetMapKeys("template") {
		generatedFiles := cnf.GetStrings("template." + templateAlias + ".generatedFiles")
		for _, generatedFileName := range generatedFiles {
			generatedFilesMap[generatedFileName] = templateAlias
		}
	}

	if *submissionFilePath != "" {
		// Specified submission file exists in generated-file list.
		if templateAlias, ok := generatedFilesMap[*submissionFilePath]; ok {
			return templateAlias
		}

		// Determine template alias to use.
		candidateAliases := make([]string, 0)
		submissionFilePathExt := filepath.Ext(*submissionFilePath)
		for _, templateAlias := range cnf.GetMapKeys("template") {
			templateFilePath := cnf.GetString("template." + templateAlias + ".codeFile")
			if submissionFilePathExt == filepath.Ext(templateFilePath) {
				candidateAliases = append(candidateAliases, templateAlias)
			}
		}

		if len(candidateAliases) == 0 {
			fmt.Println(color.RedString("error selecting solution file:"),
				"no template with extension matching '"+*submissionFilePath+"' found")
			os.Exit(1)
		}

		// Auto select when only one matching template exists.
		if len(candidateAliases) == 1 {
			return candidateAliases[0]
		}

		templateAlias := ""
		err := survey.AskOne(&survey.Select{
			Message: "Which template (alias) do you want to use?",
			Options: candidateAliases,
		}, &templateAlias)
		utils.SurveyOnInterrupt(err)

		return templateAlias
	}

	// No generated files exist.
	if len(generatedFilesMap) == 0 {
		fmt.Println(color.RedString("error selecting submission file:"),
			"file not specified; no generated files exist in local configuration")
		os.Exit(1)
	}

	// Exactly one generated file present.
	if len(generatedFilesMap) == 1 {
		for k, v := range generatedFilesMap {
			*submissionFilePath = k
			return v
		}
	}

	// Prompt user to select file from generated-files list.
	err := survey.AskOne(&survey.Select{
		Message: "Which (generated) file do you want to use?",
		Options: utils.ExtractMapKeys(generatedFilesMap),
	}, submissionFilePath)
	utils.SurveyOnInterrupt(err)

	return generatedFilesMap[*submissionFilePath]
}

func runShellScript(script string, timeout time.Duration,
	stdin io.Reader, stdout, stderr io.Writer) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmds, err := shellquote.Split(script)
	if err != nil {
		return 0, err
	}

	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// time execution of command.
	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start)

	if ctx.Err() == context.DeadlineExceeded {
		// Timeout took place.
		return elapsed, ctx.Err()
	}

	return elapsed, err
}
