package cmd

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wantguns/unfold/api"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get your account details",
	Long:  "Save your account details in the config file",
	Run:   userCmdHandler,
}

func userCmdHandler(cmd *cobra.Command, args []string) {

	user_uuid, err := api.User()
	if err != nil {
		log.Error().Err(err).Msg("Refresh response: ")
		runtime.Goexit()
	}

	viper.Set("fold_user.uuid", user_uuid)

	fmt.Println("Fetched user info")
}
