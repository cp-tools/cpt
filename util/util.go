package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"github.com/oleiade/serrure/aes"
	"github.com/spf13/viper"
)

func SurveyErr(err error) {
	if err == terminal.InterruptErr {
		fmt.Println("Interrupted")
		os.Exit(1)
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func StartCountdown(dur time.Duration, msg string) {
	writer := uilive.New()
	writer.Start()
	for ; dur.Seconds() > 0; dur -= time.Second {
		fmt.Fprintln(writer, msg, dur.String())
		time.Sleep(time.Second)
	}
	writer.Stop()
}

func ExtractMapKeys(varMap interface{}) (data []string) {
	keys := reflect.ValueOf(varMap).MapKeys()
	for _, key := range keys {
		data = append(data, key.String())
	}
	return
}

// DetectSpfr returns specifier and root workspace
func DetectSpfr(args []string) (string, string) {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) != 0 {
		return strings.Join(args, " "), currDir
	}

	tmp := strings.Split(currDir, string(os.PathSeparator))
	for c := len(tmp) - 1; c >= 0; c-- {
		if tmp[c] == "codeforces" {
			spfr, _ := DetectSpfr(tmp[c+1:])
			return spfr, strings.TrimSuffix(currDir, filepath.Join(tmp[c:]...))
		}
	}
	return "", currDir
}

func Encrypt(key string, data string) (string, error) {
	enc, err := aes.NewAES256Encrypter(key, nil)
	if err != nil {
		return "", err
	}

	ed, err := enc.Encrypt([]byte(data))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ed), err
}

func Decrypt(key string, data string) (string, error) {
	ed, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	dec := aes.NewAES256Decrypter(key)
	dd, err := dec.Decrypt(ed)
	if err != nil {
		return "", err
	}

	return string(dd), nil
}

func FindCodeFiles(file string) (string, error) {
	// file provided. Check if exists
	if len(file) != 0 {
		if _, err := os.Stat(file); err != nil {
			return "", err
		}
		return file, nil
	}

	exts := map[string]bool{}
	// store extensions of each configured template
	for _, data := range viper.GetStringMap("templates") {
		fl := data.(map[string]interface{})["file"].(string)
		exts[filepath.Ext(fl)] = true
	}

	// determine all possible files to use
	fl, err := os.Open(".")
	if err != nil {
		return "", err
	}
	defer fl.Close()

	list, err := fl.Readdirnames(0)
	if err != nil {
		return "", err
	}

	var files []string
	for _, flname := range list {
		if exts[filepath.Ext(flname)] {
			files = append(files, flname)
		}
	}

	if len(files) == 0 {
		return "", fmt.Errorf("Could not find any code files in current directory")
	} else if len(files) == 1 {
		// just one file, no selection required
		return files[0], nil
	}

	SurveyErr(survey.AskOne(&survey.Select{
		Message: "Select code file:",
		Options: files,
	}, &file))

	return file, nil
}

// FindTemplateToUse returns alias of template to use
func FindTemplateToUse(file string) (string, error) {
	var tmpltAlias string
	var tmpltsAlias []string
	for alias, data := range viper.GetStringMap("templates") {
		fl := data.(map[string]interface{})["file"].(string)
		if filepath.Ext(file) == filepath.Ext(fl) {
			tmpltsAlias = append(tmpltsAlias, alias)
		}
	}

	if len(tmpltsAlias) == 0 {
		return "", fmt.Errorf("No templates matching file %v found", file)
	} else if len(tmpltsAlias) == 1 {
		return tmpltsAlias[0], nil
	}

	SurveyErr(survey.AskOne(&survey.Select{
		Message: "Select template (alias) to use:",
		Options: tmpltsAlias,
	}, &tmpltAlias))

	return tmpltAlias, nil
}

func FindInpOutFiles(inpf, outf *[]string) {
	if len(*inpf) != 0 {
		return
	}

	files, err := filepath.Glob("*")
	if err != nil {
		panic(err)
	}

	inpRe := regexp.MustCompile(`^\d+.in$`)
	outRe := regexp.MustCompile(`^\d+.out$`)

	for _, file := range files {
		if inpRe.Match([]byte(file)) {
			*inpf = append(*inpf, file)
		} else if outRe.Match([]byte(file)) {
			*outf = append(*outf, file)
		}
	}
	return
}

func ToByte(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func SliceContains(key string, data []string) bool {
	for _, v := range data {
		if strings.EqualFold(key, v) {
			return true
		}
	}
	return false
}

func BrowserOpen(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
	return
}

func Diff(ouf, out string) string {
	t := uitable.New()
	t.Separator = " | "
	t.Wrap = true

	t.AddRow(HeaderCol("Output"), HeaderCol("Answer"))
	oufData := strings.Split(strings.TrimSpace(ouf), " ")
	outData := strings.Split(strings.TrimSpace(out), " ")

	c := 0
	for c < len(oufData) && c < len(outData) {
		t.AddRow(oufData[c], outData[c])
		c++
	}

	for c < len(oufData) {
		t.AddRow(oufData[c], "")
		c++
	}

	for c < len(outData) {
		t.AddRow("", outData[c])
		c++
	}

	return t.String()
}

func HeaderCol(data string) string {
	// simple blue bold with underline
	col := color.New(color.FgBlue).Add(color.Underline)
	return col.Sprint(data)
}
