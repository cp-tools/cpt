package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/cmd/cf"
	"github.com/cp-tools/cpt/util"
	"github.com/fatih/color"
	"github.com/infixint943/cookiejar"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cpt",
		Short: "Comprehensive, yet lightweight helper for CP!",
		Long: "Built with GO, 'cpt' is a CLI helper for most competitive coding\n" +
			"websites, that enable you to fetch sample tests, submit solutions\n" +
			"watch standings and user status, open problems page in the browser\n" +
			"all without leaving your terminal!\n\n" +
			"Automates all the boring parts of CP, so that you don't have to!\n" +
			"A little venomous advantage doesn't hurt, does it?\n\n" +
			"Built by CP'ers, built for CP'ers.",
		Version: "v0.12.0",
	}

	cfgDir string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(cf.RootCmd)
}

func initConfig() {
	// set cfgDir folder
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfgDir = filepath.Join(dir, "cpt")
	os.Mkdir(cfgDir, os.ModePerm)

	// configure default settings
	viper.SetDefault("default_template", "none")
	viper.SetDefault("gen_on_fetch", false)
	viper.SetDefault("enable_colorization", false)

	// load global settings
	cfgFile := filepath.Join(cfgDir, "cpt.json")
	viper.SetConfigFile(cfgFile)
	viper.SafeWriteConfig()
	viper.ReadInConfig()

	// set application wide colorization to use
	color.NoColor = !viper.GetBool("enable_colorization")

	// set global proxy configuration
	if viper.IsSet("proxy_url") == false {
		// use system proxy as global defaults
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	} else {
		proxyURL, _ := url.Parse(viper.GetString("proxy_url"))
		http.DefaultTransport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	// set global cookiejar to use
	jar, _ := cookiejar.New(nil)
	jar.UnmarshalJSON(util.ToByte(viper.Get("cookies")))
	http.DefaultClient.Jar = jar
	// if cookies type is not set in the beginning
	// rewriting parsed values (without changing cookies)
	// causes the cookies data to be written as a map.
	viper.Set("cookies", jar)

	// load modules of each website
	cf.InitConfig(cfgDir)
}
