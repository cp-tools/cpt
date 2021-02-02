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
	"github.com/shirou/gopsutil/process"
)

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

// Execute runs the command with resource restrictions.
func Execute(dir, command string,
	stdin io.Reader, stdout, stderr io.Writer,
	timeLimit time.Duration, memoryLimit uint64) (time.Duration, uint64, error) {

	cmds, err := shellquote.Split(command)
	if err != nil {
		return 0, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = stdin, stdout, stderr

	timer := time.Now()
	cmd.Start()

	// Actively check for MLE.
	ch := make(chan error)
	go func() { ch <- cmd.Wait() }()

	err = nil
	var memoryConsumed uint64
	for running := true; running; {
		select {
		case err = <-ch:
			running = false

		default:
			pid := int32(cmd.Process.Pid)
			if p, err := process.NewProcess(pid); err == nil {
				if m, err := p.MemoryInfo(); err == nil {
					if m.RSS > memoryConsumed {
						memoryConsumed = m.RSS
					}
				}
			}

			if memoryConsumed > memoryLimit {
				cmd.Process.Kill()
			}
		}
	}

	timeConsumed := time.Since(timer)

	if ctx.Err() == context.DeadlineExceeded {
		err = ctx.Err()
	}

	return timeConsumed, memoryConsumed, err
}
