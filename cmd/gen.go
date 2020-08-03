package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/cmd/cf"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Args:  cobra.NoArgs,
	Short: "Generate template code in current directory",
	Long: `Creates a new template file in the current directory with name <folder-name>.
If file already exists, creates file <folder-name>_<i> where 'i' iterates from 1 till the number
for which no file of the given name exists.

Usage examples:
cpt gen                    
                            Creates the default configured template
cpt gen -t fft             
                            Creates configured template of alias 'fft'
`,
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
		if tmplt, _ := lflags.GetString("template"); tmplt == "" {
			defTmplt := viper.GetString("default_template")
			if defTmplt == "" {
				return fmt.Errorf("Invalid flags - no template specified")
			}
			lflags.Lookup("template").Value.Set(defTmplt)
		}

		// check if given '--template' flag is valid
		allTmplts := util.ExtractMapKeys(viper.GetStringMap("templates"))
		if tmplt, _ := lflags.GetString("template"); !util.SliceContains(tmplt, allTmplts) {
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
	tmplt, _ := lflags.GetString("template")
	if viper.IsSet("templates."+tmplt) == false {
		fmt.Println("Template '", tmplt, "' not configured!")
		os.Exit(1)
	}
	tmpltConfig := viper.GetStringMap("templates." + tmplt)

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
			fmt.Println("File", fName, "already exists in directory")
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

			fmt.Println("Created file", fName, "in current directory")
			break
		}

		fName = fmt.Sprintf("%v_%d%v", fileBase, c, fileExt)
	}
}
