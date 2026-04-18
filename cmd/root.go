package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "git-resume",
	Short: "Git commit analyzer for resume generation",
	Long: `Git Resume Analyzer transforms your Git commit history
into impactful resume bullet points using Claude AI.

It analyzes commits, filters out noise, and generates
STAR-format achievement statements perfect for resumes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")
	rootCmd.PersistentFlags().String("repo", "", "path to git repository")
	rootCmd.PersistentFlags().String("db", "./data/cache.db", "path to SQLite database")
	rootCmd.PersistentFlags().String("output", "./output", "output directory")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose logging")

	viper.BindPFlag("repo", rootCmd.PersistentFlags().Lookup("repo"))
	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
