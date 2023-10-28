package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type TransactionsResponse struct {
	Meta struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
		URI       string    `json:"uri"`
	} `json:"meta"`
	Data struct {
		Transactions []struct {
			UUID                         string      `json:"uuid"`
			Amount                       float64     `json:"amount"`
			CurrentBalance               float64     `json:"current_balance"`
			TxnTimestamp                 time.Time   `json:"txn_timestamp"`
			TxnDate                      time.Time   `json:"txn_date"`
			IsValidTime                  bool        `json:"is_valid_time"`
			Mode                         string      `json:"mode"`
			Type                         string      `json:"type"`
			Narration                    string      `json:"narration"`
			Category                     interface{} `json:"category"`
			CategoryIcon                 interface{} `json:"category_icon"`
			Merchant                     interface{} `json:"merchant"`
			MerchantIcon                 interface{} `json:"merchant_icon"`
			MerchantAddress              interface{} `json:"merchant_address"`
			AccountID                    string      `json:"account_id"`
			Tags                         interface{} `json:"tags"`
			Kind                         string      `json:"kind"`
			FinancialInformationProvider struct {
				UUID    string `json:"uuid"`
				Name    string `json:"name"`
				FipID   string `json:"fip_id"`
				LogoURL string `json:"logo_url"`
			} `json:"financial_information_provider"`
			Notes                interface{}   `json:"notes"`
			ExcludedFromCashFlow bool          `json:"excluded_from_cash_flow"`
			IsBookmarked         bool          `json:"is_bookmarked"`
			TransactionID        string        `json:"transaction_id"`
			Reference            string        `json:"reference"`
			ExtractedTime        interface{}   `json:"extracted_time"`
			Summary              string        `json:"summary"`
			InvalidTxnID         bool          `json:"invalid_txn_id"`
			BeforeFoldAccount    bool          `json:"before_fold_account"`
			Via                  interface{}   `json:"via"`
			AccountIn            interface{}   `json:"account_in"`
			RefundStatus         string        `json:"refund_status"`
			NotifyOnRefund       bool          `json:"notify_on_refund"`
			RefundReceivedOn     interface{}   `json:"refund_received_on"`
			Receipts             []interface{} `json:"receipts"`
			GroupIds             interface{}   `json:"group_ids"`
			ContactID            interface{}   `json:"contact_id"`
		} `json:"transactions"`
		Counts []struct {
			Date              string `json:"date"`
			Total             int    `json:"total"`
			BeforeFoldAccount int    `json:"before_fold_account"`
			AfterFoldAccount  int    `json:"after_fold_account"`
		} `json:"counts"`
		Total         int         `json:"total"`
		SearchSummary interface{} `json:"search_summary"`
		After         string      `json:"after"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

type FilteredTransactions struct {
	UUID           string    `json:"uuid"`
	Amount         float64   `json:"amount"`
	CurrentBalance float64   `json:"current_balance"`
	TxnTimestamp   time.Time `json:"txn_timestamp"`
	Type           string    `json:"type"`
	Account        string    `json:"account"`
	Merchant       string    `json:"merchant"`
}

type TransactionsReturn struct {
	Transactions []FilteredTransactions
}

func randomCursor() string {
	return strconv.Itoa(10000000 + rand.Intn(89999999))
}

func filterTransactions(raw TransactionsResponse, since time.Time) []FilteredTransactions {
	transactions := make([]FilteredTransactions, 0)

	t := raw.Data.Transactions
	for i := 0; i < len(t); i++ {

		if t[i].TxnTimestamp.Before(since) {
			break
		}

		transaction := FilteredTransactions{
			UUID:           t[i].UUID,
			Amount:         t[i].Amount,
			Type:           t[i].Type,
			Account:        t[i].FinancialInformationProvider.Name,
			Merchant:       t[i].Narration,
			TxnTimestamp:   t[i].TxnTimestamp,
			CurrentBalance: t[i].CurrentBalance,
		}

		// Use Fold's F1 classifier if this transaction was classified
		if t[i].Merchant != nil {
			transaction.Merchant = t[i].Merchant.(string)
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

func Transactions(uuid string, since time.Time, till time.Time) (TransactionsReturn, error) {

	RefreshOrFail()

	req, _ := APIRequest("GET", Url("/v1/users/"+uuid+"/transactions"), nil)

	randomCursor := randomCursor() + "," + till.Format(time.RFC3339)
	log.Debug().Msg("Random cursor generated: " + randomCursor)
	randomCursorB64 := b64.StdEncoding.EncodeToString([]byte(randomCursor))
	log.Debug().Msg("Random cursor base64 generated: " + randomCursorB64)

	q := req.URL.Query()
	q.Add("filter", "all")
	q.Add("count_by", "month")
	q.Add("after", randomCursorB64)
	req.URL.RawQuery = q.Encode()

	resp, err := Client.Do(req)

	if err != nil {
		return TransactionsReturn{}, err
	} else {

		log.Debug().Msgf("Transactions response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return TransactionsReturn{}, errors.New(resp.Status)
		}

		data := TransactionsResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		log.Debug().Msgf("Transactions response body: %+v", data.Data.Transactions[0].TxnTimestamp)

		var ret TransactionsReturn
		ret.Transactions = make([]FilteredTransactions, 0)
		ret.Transactions = append(ret.Transactions, filterTransactions(data, since)...)

		for data.Data.Transactions[len(data.Data.Transactions)-1].TxnTimestamp.After(since) {
			log.Debug().Msg("Fetching older transactions")

			log.Debug().Msg("New cursor base64: " + data.Data.After)
			q.Set("after", data.Data.After)
			req.URL.RawQuery = q.Encode()

			resp, err := Client.Do(req)
			log.Debug().Msgf("Transactions response status: %+v", resp.StatusCode)

			if err != nil {
				log.Warn().Msg("Failed to fetch older transactions")
				break
			} else {
				if resp.StatusCode/100 != 2 {
					log.Warn().Msgf("Failed to fetch older transactions, status code: %+v", resp.StatusCode)
					break
				}

				json.NewDecoder(resp.Body).Decode(&data)
				log.Debug().Msgf("Transactions response body: %+v", data.Data.Transactions[0].TxnTimestamp)

				ret.Transactions = append(ret.Transactions, filterTransactions(data, since)...)
			}
		}

		return ret, nil
	}
}
