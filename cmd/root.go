package cmd

import (
	"github.com/fzdwx/gh-sp/api"
	"github.com/fzdwx/gh-sp/cmd/repo"
	"github.com/fzdwx/gh-sp/cmd/search"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "code-github-workspace",
	Short: "Github workspace is an interactive github cli program",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(search.New())
	rootCmd.AddCommand(repo.New())

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.code-github-workspace.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "use debug log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	f, err := os.OpenFile("code.log", os.O_RDWR|os.O_CREATE, 0777)
	cobra.CheckErr(err)
	api.InitLog(debug, f)
	api.InitBrowser(f)
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := os.UserHomeDir()
	//
	//	cobra.CheckErr(err)
	//
	//	// Search config in home directory with name ".code-github-workspace" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigType("yaml")
	//	viper.SetConfigName(".code-github-workspace")
	//}
	//
	//viper.AutomaticEnv() // read in environment variables that match
	//
	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err != nil {
	//	panic(err)
	//}
	//
	//var cfg config.Config
	//if err := viper.Unmarshal(&cfg); err != nil {
	//	panic(err)
	//}
	//
	//config.Init(&cfg)
	//
	//log.Debug().Str("config file", viper.ConfigFileUsed()).Msg("load config success")
}
