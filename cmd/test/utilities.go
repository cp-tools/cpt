package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cp-tools/cpt/packages/conf"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
)

// extractGeneratedFiles returns a map consisting of previously generated
// files in local configuration, along with the corresponding template alias.
func extractGeneratedFiles(cnf *conf.Conf) map[string]string {
	data := make(map[string]string)
	// Extract all template aliases.
	for _, alias := range cnf.GetMapKeys("template") {
		// Extract and all generated files of alias.
		generatedFiles := cnf.GetMapKeys("template." + alias + ".generatedFiles")
		for _, fileName := range generatedFiles {
			data[fileName] = alias
		}
	}
	return data
}

// SelectCodeFile returns file name and template of code file to use, based
// on configured templates and passed 'filePath' value.
// 'filePath' must point to a valid file.
func SelectCodeFile(filePath string, cnf *conf.Conf) (fileName string, alias string) {
	// Find all generated code files in local configurations.
	generatedFilesMap := extractGeneratedFiles(cnf)
	// Check if filePath exists in generatedFilesMap.
	if _, ok := generatedFilesMap[filePath]; !ok {
		// Try to auto select code file, if not specified.
		if filePath == "" {
			if len(generatedFilesMap) != 1 {
				fmt.Println(color.RedString("error selecting solution file:"),
					"file not specified, unable to auto-select code file from local configurations")
				os.Exit(1)
			}
			// Auto select code file to use.
			for k, v := range generatedFilesMap {
				fileName, alias = k, v
			}
			return
		}

		// Find all templates with extension matching filePath.
		aliasData := make([]string, 0)
		fileExtension := filepath.Ext(filePath)
		for _, alias := range cnf.GetMapKeys("template") {
			codeFile := cnf.GetString("template." + alias + ".codeFile")
			if fileExtension == filepath.Ext(codeFile) {
				aliasData = append(aliasData, alias)
			}
		}

		if len(aliasData) == 0 {
			fmt.Println(color.RedString("error selecting solution file:"),
				"no template with code file matching '"+filePath+"' found")
			os.Exit(1)
		} else if len(aliasData) == 1 {
			// Auto set template configuration to use.
			fileName, alias = filePath, aliasData[0]
			return
		}

		fileName = filePath
		// Prompt user to select template alias to use.
		survey.AskOne(&survey.Select{
			Message: "Which template (alias) do you want to use?",
			Options: aliasData,
		}, &alias)
		return
	}
	fileName, alias = filePath, generatedFilesMap[filePath]
	return
}

func runShellScript(script string, timeout time.Duration,
	stdin, stdout, stderr io.ReadWriter) (time.Duration, error) {
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

	return elapsed, err
}
