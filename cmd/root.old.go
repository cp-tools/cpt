package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cp-tools/cpt/cmd/cf"
	"github.com/cp-tools/cpt/util"
	"github.com/infixint943/cookiejar"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgDir string
)



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
