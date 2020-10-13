package generate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
func Generate(templateMap map[string]interface{}) {
	// Read template codeFile to variable.
	templateData, err := ioutil.ReadFile(templateMap["codeFile"].(string))
	if err != nil {
		color.Red("error reading template code file: %v", err)
		os.Exit(1)
	}

	// Extract base name of the current directory.
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unexpected error occurred: %v", err)
	}

	baseFileName := filepath.Base(currentDir)
	fileExtension := filepath.Ext(templateMap["codeFile"].(string))
	fileName := decideFileName(baseFileName, fileExtension)

	file, err := os.Create(fileName)
	if err != nil {
		color.Red("error creating code file: %v", err)
		os.Exit(1)
	}

	// Write templateData to created file.
	if _, err := file.Write(templateData); err != nil {
		color.Red("error writing to code file: %v", err)
		os.Exit(1)
	}

	color.Green("created code file: %v", fileName)
}

func decideFileName(baseFileName, fileExtension string) string {
	for fileName, i := baseFileName, 0; true; i++ {
		fullName := fileName + fileExtension
		if file, err := os.Stat(fullName); os.IsNotExist(err) || file.IsDir() {
			return fullName
		}
		fileName = fmt.Sprintf("%v_%d", baseFileName, i)
	}
	// Impossible case
	return ""
}
