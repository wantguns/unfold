package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
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
	Run:   setupTransactionsCmdHandler,
}

func init() {
	now := time.Now()
	today := now.AddDate(0, 0, 1).Format(time.DateOnly)
	yesterday := now.AddDate(0, -1, 0).Format(time.DateOnly)

	TransactionsCmd.Flags().StringP("till", "t", today, "fetch transactions till in this format: YYYY-MM-DD")
	TransactionsCmd.Flags().StringP("since", "s", yesterday, "fetch transactions since in this format: YYYY-MM-DD")
	TransactionsCmd.Flags().BoolP("db", "d", false, "Save the results in a sqlite db")
	TransactionsCmd.Flags().StringP("db-path", "D", "db.sqlite", "Sets path for the database")
	TransactionsCmd.Flags().StringP("watch", "w", "", "Set an internal cron job to trigger this command. You can use non-standard cron expressions like '@every 6h'. This will disable plaintext mode, so add a '-d' flag if you want to write to db")
}

func setupTransactionsCmdHandler(cmd *cobra.Command, args []string) {
	watch, _ := cmd.Flags().GetString("watch")
	if watch == "" {
		// Fetch transactions in a oneshot manner

		// Update the `plaintext` value
		oldArgs := os.Args[1:]
		log.Debug().Msgf("Old arguments %+v", oldArgs)
		cmd.SetArgs(append(oldArgs, "--no-plaintext"))

		transactionsCmdHandler(cmd, args)
	} else {
		log.Info().Msg("Cron job set for fetching transactions, going into daemon mode")

		// Fetch transactions once before going into cron land
		transactionsCmdHandler(cmd, args)
		till, _ := cmd.Flags().GetString("till")
		log.Info().Msgf("Fetched transactions till %s", till)

		c := cron.New()
		c.AddFunc(watch, func() {

			// Update the `till` and `plaintext` value
			now := time.Now().Format(time.DateOnly)
			oldArgs := os.Args[1:]
			log.Debug().Msgf("Old arguments %+v", oldArgs)
			cmd.SetArgs(append(oldArgs, fmt.Sprintf("--till=%s", now), "--no-plaintext"))

			transactionsCmdHandler(cmd, args)
			log.Info().Msgf("Fetched transactions till %s", now)
		})

		go c.Start()
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
	}
}

func printTransactions(t api.FilteredTransactions) {
	fmt.Printf(
		"%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		t.UUID,
		t.TxnTimestamp,
		t.Amount,
		t.Type,
		t.Merchant,
		t.CurrentBalance,
		t.Account,
	)
}

func writeToDb(t api.FilteredTransactions) {
	db.Conn.Clauses(clause.OnConflict{UpdateAll: true}).Create(&db.Transactions{
		UUID:           t.UUID,
		Timestamp:      t.TxnTimestamp,
		Amount:         t.Amount,
		Type:           t.Type,
		Merchant:       t.Merchant,
		CurrentBalance: t.CurrentBalance,
		Account:        t.Account,
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

		// If plaintext is enabled
		if noPlaintext := cmd.Flags().Lookup("no-plaintext"); noPlaintext != nil {
			printTransactions(t[i])
		}
	}
}
