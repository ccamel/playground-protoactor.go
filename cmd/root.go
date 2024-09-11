//nolint:gochecknoglobals,gochecknoinits // common pattern when using cobra library
package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "playground-protoactor.go",
	Short: "Playground for playing with protoactor (the next gen Actor Model framework)",
	Long: `Playground for playing with protoactor (the next gen Actor Model framework) in go,
	following DDD, Event Sourcing & CQRS paradigms.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) //nolint:forbidigo // common pattern when using cobra library
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.playground-protoactor.go.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err) //nolint:forbidigo // common pattern when using cobra library
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".playground-protoactor.go")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed()) //nolint:forbidigo // common pattern when using cobra library
	}
}
