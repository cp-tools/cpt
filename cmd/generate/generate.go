package generate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/pkg/conf"

	"github.com/fatih/color"
)

// Generate creates a new file in the current directory,
// using given template data.
//
// It names the file the same name as the current folder.
// However, if a file with the same name exists, it appends _<i>,
// where <i> is the smallest non-negative number for which no file exists.
//
// Spaces in the file name are also replaced with underscore.
//
// Ensure 'cnf' is local (folder) configuration.
func Generate(alias string, cnf *conf.Conf, dataMap map[string]interface{}) {
	// Extract templateMap from conf.
	templateMap, ok := cnf.Get("template." + alias).(map[string]interface{})
	if !ok {
		fmt.Println(color.RedString("unexpected error occurred:"),
			"template '"+alias+"' is not of type map[string]interface{}")
		os.Exit(1)
	}

	// Read template codeFile to variable.
	codeFile, ok := templateMap["codeFile"].(string)
	if !ok {
		fmt.Println(color.RedString("unexpected error occurred:"),
			"field codefile is not of type string")
		os.Exit(1)
	}

	templateData, err := ioutil.ReadFile(codeFile)
	if err != nil {
		color.Red("error reading template code file: %v", err)
		os.Exit(1)
	}
	templateData = updatePlaceholders(templateData, dataMap)

	// Extract base name of the current directory.
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unexpected error occurred: %v", err)
	}

	baseFileName := filepath.Base(currentDir)
	fileExtension := filepath.Ext(codeFile)
	fileName := DecideFileName(baseFileName, fileExtension)

	file, err := os.Create(fileName)
	if err != nil {
		color.Red("error creating code file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write templateData to created file.
	if _, err := file.Write(templateData); err != nil {
		color.Red("error writing to code file: %v", err)
		os.Exit(1)
	}
	// Write generated file details to local conf.
	generatedFiles := cnf.GetStrings("template." + alias + ".generatedFiles")
	generatedFiles = append(generatedFiles, fileName)

	cnf.Set("template."+alias+".generatedFiles", generatedFiles)
	cnf.WriteFile()

	color.Green("created code file: %v", fileName)
}
