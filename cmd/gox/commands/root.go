package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/gox-framework/gox/pkg/version"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gox",
	Short: "GOX Framework - Modern web development with Go and HTMX",
	Long: `GOX Framework is a modern web development framework that combines
the power of Go with the simplicity of HTMX for ultra-fast Server-Side Rendered applications.

Inspired by Vue.js SFC, React and Svelte, but designed to be "easy to learn, difficult to master".

Examples:
  gox new project my-app        Create a new project
  gox new module user-mgmt      Create a new module  
  gox generate page dashboard   Generate a new page
  gox dev                       Start development server
  gox build                     Build for production`,
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is gox.config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	
	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	
	// Version command
	rootCmd.SetVersionTemplate(version.String() + "\n")
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("gox.config")
	}

	// Environment variables
	viper.SetEnvPrefix("GOX")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}