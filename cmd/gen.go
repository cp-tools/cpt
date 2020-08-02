package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/cmd/cf"
	"github.com/cp-tools/cpt/util"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate template code in current directory",
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(genCmd)
	// add flags here
	genCmd.Flags().StringP("template", "t", "", "Template (by alias) to use")
	genCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command,
		_ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		allTmplts := util.ExtractMapKeys(viper.GetStringMap("templates"))
		return allTmplts, cobra.ShellCompDirectiveNoFileComp
	})

	genCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lflags := genCmd.Flags()
		// set template to default template if flag '--template' not set
		if lflags.MustGetString("template") == "" {
			defTmplt := viper.GetString("default_template")
			if defTmplt == "" {
				return fmt.Errorf("Invalid flags - no template specified")
			}
			lflags.Lookup("template").Value.Set(defTmplt)
		}

		// check if given '--template' flag is valid
		allTmplts := util.ExtractMapKeys(viper.GetStringMap("templates"))
		if tmplt := lflags.MustGetString("template"); !util.SliceContains(tmplt, allTmplts) {
			return fmt.Errorf("Invalid flags - template '%v' not found", tmplt)
		}

		gen(lflags)
		return nil
	}

	// pass gen to subcommands
	cf.GenFunc = gen
}

func gen(lflags *pflag.FlagSet) {
	// get template configuration to use
	tmpltConfig := viper.GetStringMap("templates." + lflags.MustGetString("template"))

	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// generate file with foldername+index.extension
	fileBase := filepath.Base(currDir)
	fileExt := filepath.Ext(tmpltConfig["file"].(string))
	for fName, c := fileBase+fileExt, 1; true; c++ {

		if _, err := os.Stat(fName); os.IsNotExist(err) == false {
			color.Yellow("File %v exists in directory", fName)
		} else {
			data, err := ioutil.ReadFile(tmpltConfig["file"].(string))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = ioutil.WriteFile(fName, data, 0644)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			color.Green("Created template file %v", fName)
			break
		}

		fName = fmt.Sprintf("%v_%d%v", fileBase, c, fileExt)
	}
}
