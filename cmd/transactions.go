package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"

	"github.com/wantguns/unfold/api"
	"github.com/wantguns/unfold/db"
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
	TransactionsCmd.Flags().BoolP("db", "d", false, "Save the results in a sqlite db")
	TransactionsCmd.Flags().StringP("db-path", "D", "db.sqlite", "Sets path for the database")
}

func printTransactions(t api.FilteredTransactions) {
	fmt.Printf(
		"%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		t.UUID,
		t.TxnTimestamp,
		t.Amount,
		t.Type,
		t.Sender,
		t.CurrentBalance,
		t.Receiver,
	)
}

func writeToDb(t api.FilteredTransactions) {
	db.Conn.Clauses(clause.OnConflict{UpdateAll: true}).Create(&db.Transactions{
		UUID:           t.UUID,
		Timestamp:      t.TxnTimestamp,
		Amount:         t.Amount,
		Type:           t.Type,
		Sender:         t.Sender,
		CurrentBalance: t.CurrentBalance,
		Receiver:       t.Receiver,
	})
}

func transactionsCmdHandler(cmd *cobra.Command, args []string) {

	uuid := viper.GetString("fold_user.uuid")

	// till Flag
	tillStr, _ := cmd.Flags().GetString("till")
	till, err := time.Parse(time.DateOnly, tillStr)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid time format `till`: %+v", tillStr)
		runtime.Goexit()
	}
	if till.After(time.Now()) {
		till = time.Now()
	}

	// since Flag
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

	// db Flag
	writeDb, _ := cmd.Flags().GetBool("db")
	dbPath, _ := cmd.Flags().GetString("db-path")
	if writeDb {
		db.Init(dbPath)
		log.Debug().Msgf("Database path %s", dbPath)
	}

	transactions, err := api.Transactions(uuid, since, till)
	if err != nil {
		log.Error().Err(err).Msg("Refresh response: ")
		runtime.Goexit()
	}

	t := transactions.Transactions
	for i := 0; i < len(t); i++ {
		// Insert into db
		if writeDb {
			writeToDb(t[i])
		}

		// always printTransactions
		printTransactions(t[i])
	}
}
