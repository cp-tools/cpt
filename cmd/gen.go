package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/cmd/cf"
	"github.com/cp-tools/cpt/util"

	"github.com/spf13/cobra"
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

	var tmpltFlag string
	genCmd.Flags().StringVarP(&tmpltFlag, "template", "t", "",
		"Select template configuration (by alias name) to use")

	genCmd.Run = func(cmd *cobra.Command, args []string) {
		gen(tmpltFlag)
	}

	// pass gen to subcommands
	cf.GenFunc = gen
}

func gen(tmplt string) {
	if len(tmplt) == 0 {
		tmplt = viper.GetString("default_template")
	}

	allTmplts := util.ExtractMapKeys(viper.GetStringMap("templates"))
	if !util.SliceContains(tmplt, allTmplts) {
		fmt.Println("Select a valid template [alias] to generate")
		os.Exit(2)
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
