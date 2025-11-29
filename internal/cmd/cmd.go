package cmd

import (
	"log"

	"github.com/innovationmech/simple-cli/internal/cmd/serve"
	"github.com/innovationmech/simple-cli/internal/cmd/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func NewRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "simple-cli",
		Short: "Simple CLI",
		Long:  "Simple CLI",
	}

	rootCmd.PersistentFlags().Int("port", 8080, "Port to listen on")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.yaml", "config file (default is ./config.yaml)")

	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	viper.SetEnvPrefix("SIMPLE_CLI")
	viper.AutomaticEnv()

	rootCmd.AddCommand(version.NewVersionCmd())
	rootCmd.AddCommand(serve.NewServeCmd())

	return rootCmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}
