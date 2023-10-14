package api

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"time"
)

type AvailabilityResponse struct {
	Meta struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
		URI       string    `json:"uri"`
	} `json:"meta"`
	Data struct {
		Accounts []struct {
			AccountID                 string    `json:"account_id"`
			IsHidden                  bool      `json:"is_hidden"`
			TransactionAvailableSince time.Time `json:"transaction_available_since"`
			TransactionAvailableTill  time.Time `json:"transaction_available_till"`
			Version                   int       `json:"version"`
		} `json:"accounts"`
		Version int `json:"version"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

func Availability(uuid string) (time.Time, time.Time, error) {

	RefreshOrFail()

	log.Debug().Msg("Availability URL: " + Url("/v1/users/"+uuid+"/transactions/available_since"))

	req, _ := APIRequest("GET", Url("/v1/users/"+uuid+"/transactions/available_since"), nil)
	resp, err := Client.Do(req)

	if err != nil {
		return time.Time{}, time.Time{}, err
	} else {

		log.Debug().Msgf("Availability response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return time.Time{}, time.Time{}, errors.New(resp.Status)
		}
		defer resp.Body.Close()

		data := AvailabilityResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		log.Debug().Msgf("Availability response: %+v", data)

		minTransactionAvailableSince := time.Now()
		maxTransactionAvailableTill := time.Time{}
		accounts := data.Data.Accounts

		for i := 0; i < len(accounts); i++ {
			if accounts[i].TransactionAvailableSince.Before(minTransactionAvailableSince) {
				minTransactionAvailableSince = accounts[i].TransactionAvailableSince
			}

			if accounts[i].TransactionAvailableTill.After(maxTransactionAvailableTill) {
				maxTransactionAvailableTill = accounts[i].TransactionAvailableTill
			}
		}

		return minTransactionAvailableSince, maxTransactionAvailableTill, nil
	}
}
