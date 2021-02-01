package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/pkg/conf"
	"github.com/cp-tools/cpt/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
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
