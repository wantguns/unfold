package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/wantguns/unfold/cmd"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Error().Msg("SIGINT recieved, shutting down")
			viper.WriteConfig()
			os.Exit(1)
		}
	}()

	defer viper.WriteConfig()

	cmd.Execute()
}
