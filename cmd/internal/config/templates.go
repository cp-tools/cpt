package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cp-tools/cpt/util"
	"github.com/fatih/color"
	"github.com/kballard/go-shellquote"
	"github.com/mitchellh/go-homedir"
)

// AddTemplate inserts a new template into the template map.
func AddTemplate(dataMap map[string]interface{}) {
	// template alias - what template will be called
	alias := ""
	survey.AskOne(&survey.Input{
		Message: "What do you want to call the template?",
		Help: `Set the alias you wish to refer the template by.
It must be unique, and should not contain whitespaces.`,
	}, &alias, survey.WithValidator(func(ans interface{}) error {
		if ans.(string) == "" {
			return fmt.Errorf("alias value is required")
		}
		if strings.Contains(ans.(string), " ") {
			return fmt.Errorf("alias should not have whitespace")
		}

		for val := range dataMap {
			if ans.(string) == val {
				return fmt.Errorf("template with given alias exists")
			}
		}
		return nil
	}))

	// select file containing template code
	template := make(map[string]interface{})
	survey.Ask([]*survey.Question{
		{
			Name: "codeFile",
			Prompt: &survey.Input{
				Message: "Which code file should be used as the template?",
				Help:    `The code present in this file will be used as the template.`,
			},
			Validate: func(ans interface{}) error {
				templateCodePath, err := homedir.Abs(ans.(string))
				if err != nil {
					return err
				}
				// check if path corresponds to a valid file
				if _, err := os.Open(templateCodePath); err != nil {
					return err
				}

				return nil
			},
			Transform: func(ans interface{}) (newAns interface{}) {
				templateCodePath, _ := homedir.Abs(ans.(string))
				return templateCodePath
			},
		},
		{
			Name: "preScript",
			Prompt: &survey.Input{
				Message: "What script should be run, before testing the solution code?",
				Help: `
When testing your solution using 'cpt test', the following takes place:

- The prescript (current prompt) is run at the start, exactly once.
This is mostly used for compiling your solution file, before testing.
Intepreted languages (like python) can omit this.

Simple compilation prescripts that you can use (without quotes) are:
C++ ==> 'g++ -std=c++11 -Wall {{.file}}'
Java ==> 'javac {{.file}}'

- The runscript (next prompt) is run per test case.
See explanation in help section of next prompt.

- The postscript (final prompt) is run at the end, exactly once.
See explanation in help section of final prompt.`,
			},
			Validate: func(ans interface{}) error {
				// check if script is well formed
				_, err := shellquote.Split(ans.(string))
				return err
			},
		},
		{
			Name: "runScript",
			Prompt: &survey.Input{
				Message: "what script should be run, to execute the solution code?",
				Help: `
When testing your solution using 'cpt test', the following takes place:

- The prescript (previous prompt) is run at the start, exactly once.
See explanation in help section of previous prompt.

- The runscript (current prompt) is run per test case.
This is the command used for running your code/executable.
For each sample test, this script is run once, and the verdict of
your solution for the input is determined.
This is a required field.

Example runscripts that you can use (without quotes) are:
Python ==> 'python3 {{.file}}'
Java ==> 'java {{.fileNoExt}}'
C++ ==> './a.out' (linux) or 'a.exe' (windows)

Note that, you need not redirect input/output in your code
and instead read and write to stdin and stdout respectively.

- The postscript (final prompt) is run at the end, exactly once.
See explanation in help section of final prompt.`,
			},
			Validate: func(ans interface{}) error {
				if ans.(string) == "" {
					return fmt.Errorf("runscript value is required")
				}

				_, err := shellquote.Split(ans.(string))
				return err
			},
		},
		{
			Name: "postscript",
			Prompt: &survey.Input{
				Message: "What script should be run, after testing the solution code?",
				Help: `
When testing your solution using 'cpt test', the following takes place:

- The prescript (first prompt) is run at the start, exactly once.
See explanation in help section of first prompt.

- The runscript (previous prompt) is run per test case.
See explanation in help section of previous prompt.

- The postscript (current prompt) is run at the end, exactly once.
Generally, this is used to clean up residual files like executables, if any.
Intepreted languages (like python) can omit this.

Example postscripts you can use (without quotes) are:
C++ ==> 'rm a.out' (linux) or 'rem a.exe' (windows)
Java ==> 'rm {{.fileNoExt}}' (linux)
		 'rem {{.fileNoExt}}' (windows)
`,
			},
			Validate: func(ans interface{}) error {
				// check if script is well formed
				_, err := shellquote.Split(ans.(string))
				return err
			},
		},
	}, &template)
	// write newTemplate to all templates
	dataMap[alias] = template
	return
}

// RemoveTemplate deletes selected template from template map.
func RemoveTemplate(dataMap map[string]interface{}) {
	if len(dataMap) == 0 {
		color.Red("no templates present")
		os.Exit(1)
	}

	alias := ""
	survey.AskOne(&survey.Select{
		Message: "Which template do you want to delete?",
		Options: util.ExtractMapKeys(dataMap),
	}, &alias)
	// remove alias named template
	delete(dataMap, alias)
	return
}
