package generate

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"
)

func DecideFileName(baseFileName, fileExtension string) string {
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

func updatePlaceholders(str []byte, dataMap map[string]interface{}) []byte {
	// Add generic placeholders.
	dataMap["date"] = time.Now().Format("02.01.2006")
	dataMap["time"] = time.Now().Format("15:04")

	var out bytes.Buffer
	tmplt := template.New("")
	tmplt.Parse(string(str))
	tmplt.Execute(&out, dataMap)

	return out.Bytes()
}