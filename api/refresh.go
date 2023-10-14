package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type RefreshPayload struct {
	Error
	Refresh string `json:"refresh_token,omitempty"`
}

type RefreshResponse struct {
	Error
	Meta struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
		URI       string    `json:"uri"`
	} `json:"meta"`
	Data struct {
		TokenType    string    `json:"token_type"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		ExpiresAt    time.Time `json:"expires_at"`
	} `json:"data"`
}

func Refresh(refresh string) (string, string, error) {

	refreshPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(refreshPayloadBuf).Encode(&RefreshPayload{Refresh: refresh})

	req, _ := APIRequest("POST", Url("/v1/auth/tokens/refresh"), refreshPayloadBuf)
	resp, err := Client.Do(req)

	if err != nil {
		return "", "", err
	} else {

		log.Debug().Msgf("Refresh response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return "", "", errors.New(resp.Status)
		}

		data := RefreshResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		log.Debug().Msgf("Verify OTP response: %+v", data)
		return data.Data.AccessToken, data.Data.RefreshToken, nil
	}
}

func RefreshOrFail() {
	refresh_token := viper.GetString("token.refresh")

	access, refresh, err := Refresh(refresh_token)
	if err != nil {
		log.Error().Err(err).Msg("Failed at refreshing token")
		runtime.Goexit()
	}

	viper.Set("token.access", access)
	viper.Set("token.refresh", refresh)

	log.Debug().Msg("Refreshed tokens")
}
