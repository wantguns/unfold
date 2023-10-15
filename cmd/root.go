package cmd

import (
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	enableDebug bool

	rootCmd = &cobra.Command{
		Use:   "unfold",
		Short: "An unofficial cli client for fold.money",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/unfold/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&enableDebug, "debug", "v", os.Getenv("DEBUG") == "true", "Enable debug mode")
	rootCmd.AddCommand(LoginCmd, RefreshCmd, UserCmd, AvailabilityCmd, TransactionsCmd)
}

func initConfig() {

	// Debug Flag
	if enableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Config File
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgDir, err := os.UserConfigDir()
		cobra.CheckErr(err)
		dir := cfgDir + "/unfold"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create the config directory")
				runtime.Goexit()
			}
		}

		viper.AddConfigPath(dir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		viper.SetDefault("device_hash", uuid.NewString())
		viper.SafeWriteConfig()
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
