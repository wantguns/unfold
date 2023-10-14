package cmd

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wantguns/unfold/api"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your fold account",
	Long:  `This command must be run before running any other command to authenticate the CLI`,
	Run:   loginCmdHandler,
}

func loginCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Print("Enter the phone number associated with your fold account: ")
	phone := bufio.NewScanner(os.Stdin)
	phone.Scan()

	err := api.Login("+91" + phone.Text())
	if err != nil {
		log.Error().Err(err).Msg("Login response: ")
		runtime.Goexit()
	}

	fmt.Print("Login request successful, enter OTP: ")
	otp := bufio.NewScanner(os.Stdin)
	otp.Scan()

	access, refresh, err := api.VerifyOtp("+91"+phone.Text(), otp.Text())
	if err != nil {
		log.Error().Err(err).Msg("Verify otp response: ")
		runtime.Goexit()
	}

	viper.Set("token.access", access)
	viper.Set("token.refresh", refresh)

	log.Debug().Msg("Fetching user info")
	user_uuid, err := api.User()
	if err != nil {
		log.Error().Err(err).Msg("Refresh response: ")
		runtime.Goexit()
	}

	viper.Set("fold_user.uuid", user_uuid)

	fmt.Println("Login successful !")
}
