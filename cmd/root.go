package cmd

import (
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "unfold",
		Short: "An unofficial cli client for fold.money",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/unfold/config.yaml)")
	rootCmd.AddCommand(LoginCmd, RefreshCmd, UserCmd, AvailabilityCmd, TransactionsCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		err := viper.ReadInConfig()
		if err != nil {
			log.Error().Err(err).Msg("Failed to read config file")
			runtime.Goexit()
		}
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

		err = viper.ReadInConfig()
		if err != nil {
			log.Error().Err(err).Msg("Failed to read config file")
			runtime.Goexit()
		}

		viper.SafeWriteConfig()
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()

}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
