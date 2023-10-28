package ledger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func WriteToFile(file string, postings []Posting) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to open ledger file %+v", file)
		runtime.Goexit()
	}

	newFetchStr := fmt.Sprintf("\n\n;\n;\tunfold_meta_fetch\n;\ttime: %s\n;\n\n\n", time.Now().String())

	if _, err := f.WriteString(newFetchStr); err != nil {
		log.Error().Err(err).Msgf("Failed to write %+v to %+v ledger file", newFetchStr, file)
		runtime.Goexit()
	}

	for _, p := range postings {

		// ; unfold_meta_start             7e15319f-144f-4f89-a6dc-8f94bfed2819
		// 2014/03/03 Internet
		//     Expenses:Utilities          1500 INR
		//     Assets:Checking
		// ; unfold_meta_end               7e15319f-144f-4f89-a6dc-8f94bfed2819
		data :=
			fmt.Sprintf("\n;\tunfold_meta_start\t%s\n%s\t%s\n\tExpenses:%s\t\t%sINR\n\tAssets:Checking:%s\n;\tunfold_meta_end\t%s\n",
				p.UUID, p.Date, p.Description, p.Merchant, p.Amount, strings.Split(p.Account, " ")[0], p.UUID)

		if _, err := f.WriteString(data); err != nil {
			log.Error().Err(err).Msgf("Failed to write %+v to %+v ledger file", data, file)
			runtime.Goexit()
		}
	}

	if err := f.Close(); err != nil {
		log.Error().Err(err).Msgf("Failed to close ledger file %+v", file)
		runtime.Goexit()
	}

}
