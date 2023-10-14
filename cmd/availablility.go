package cmd

import (
	// "fmt"
	"fmt"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wantguns/unfold/api"
)

var AvailabilityCmd = &cobra.Command{
	Use:   "availability",
	Short: "Returns a range of dates for when your banking data is available",
	Run:   availabilityCmdHandler,
}

func availabilityCmdHandler(cmd *cobra.Command, args []string) {
	uuid := viper.GetString("fold_user.uuid")
	log.Debug().Msg("User UUID: " + uuid)

	availableSince, availableTill, err := api.Availability(uuid)
	if err != nil {
		log.Error().Err(err).Msg("Available response: ")
		runtime.Goexit()
	}

	fmt.Println("Transactions available since", availableSince.Format(time.RFC822Z))
	fmt.Println("Transactions available till", availableTill.Format(time.RFC822Z))
}
