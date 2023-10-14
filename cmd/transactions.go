package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wantguns/unfold/api"
)

var TransactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Prints the transactions from all of your accounts (default period: 1 month)",
	Run:   transactionsCmdHandler,
}

func init() {
	now := time.Now()
	today := now.AddDate(0, 0, 1).Format(time.DateOnly)
	yesterday := now.AddDate(0, -1, 0).Format(time.DateOnly)

	TransactionsCmd.Flags().StringP("till", "t", today, "fetch transactions till in this format: YYYY-MM-DD")
	TransactionsCmd.Flags().StringP("since", "s", yesterday, "fetch transactions since in this format: YYYY-MM-DD")
}

func transactionsCmdHandler(cmd *cobra.Command, args []string) {

	uuid := viper.GetString("fold_user.uuid")

	tillStr, _ := cmd.Flags().GetString("till")
	till, err := time.Parse(time.DateOnly, tillStr)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid time format `till`: %+v", tillStr)
		runtime.Goexit()
	}
	if till.After(time.Now()) {
		till = time.Now()
	}

	minSince, _, err := api.Availability(uuid)
	if err != nil {
		log.Error().Err(err).Msg("Fetch Availability: ")
		runtime.Goexit()
	}
	sinceStr, _ := cmd.Flags().GetString("since")
	since, err := time.Parse(time.DateOnly, sinceStr)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid time format `since`: %+v", sinceStr)
		runtime.Goexit()
	}
	if since.Before(minSince) {
		since = minSince
	}

	transactions, err := api.Transactions(uuid, since, till)
	if err != nil {
		log.Error().Err(err).Msg("Refresh response: ")
		runtime.Goexit()
	}

	fmt.Println("Fetched transactions")

	t := transactions.Transactions
	for i := 0; i < len(t); i++ {
		fmt.Printf(
			"%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			t[i].UUID,
			t[i].TxnDate,
			t[i].TxnTimestamp,
			t[i].Amount,
			t[i].Type,
			t[i].Sender,
			t[i].CurrentBalance,
			t[i].Receiver,
		)
	}
}
