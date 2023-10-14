package cmd

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wantguns/unfold/api"
)

var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh your auth tokens",
	Run:   refreshCmdHandler,
}

func refreshCmdHandler(cmd *cobra.Command, args []string) {
	refresh_token := viper.GetString("token.refresh")

	access, refresh, err := api.Refresh(refresh_token)
	if err != nil {
		log.Error().Err(err).Msg("Refresh response: ")
		runtime.Goexit()
	}

	viper.Set("token.access", access)
	viper.Set("token.refresh", refresh)
	// viper.WriteConfig()

	fmt.Println("Refreshed auth tokens !")
}
