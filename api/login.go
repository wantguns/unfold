package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type Error struct {
	Error string `json:"error,omitempty"`
}

type LoginPayload struct {
	Error
	Phone   string `json:"phone,omitempty"`
	Channel string `json:"channel,omitempty"`
}

func Login(phone string) error {
	loginPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(loginPayloadBuf).Encode(&LoginPayload{Phone: phone, Channel: "sms"})

	req, _ := APIRequest("POST", Url("/v1/auth/otp"), loginPayloadBuf)
	resp, err := Client.Do(req)

	if err != nil {
		return err
	} else {
		log.Debug().Msgf("Login response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return errors.New(resp.Status)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Debug().Msg("Login response: " + buf.String())

		return nil
	}
}

type VerifyOtpPayload struct {
	Error
	Phone string `json:"phone,omitempty"`
	Otp   string `json:"otp,omitempty"`
}

type VerifyOtpResponse struct {
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
		NewUser      bool      `json:"new_user"`
		UserID       string    `json:"user_id"`
		UserMeta     struct {
			UUID          string      `json:"uuid"`
			FirstName     string      `json:"first_name"`
			MiddleName    interface{} `json:"middle_name"`
			LastName      string      `json:"last_name"`
			Email         string      `json:"email"`
			EmailVerified bool        `json:"email_verified"`
			GoogleLinked  bool        `json:"google_linked"`
			AppleLinked   bool        `json:"apple_linked"`
			Role          string      `json:"role"`
		} `json:"user_meta"`
	} `json:"data"`
}

func VerifyOtp(phone string, otp string) (string, string, error) {
	otpVerifyPayloadBuf := new(bytes.Buffer)
	json.NewEncoder(otpVerifyPayloadBuf).Encode(&VerifyOtpPayload{Phone: phone, Otp: otp})

	req, _ := APIRequest("POST", Url("/v1/auth/otp/verify"), otpVerifyPayloadBuf)
	resp, err := Client.Do(req)

	if err != nil {
		return "", "", err
	} else {

		log.Debug().Msgf("Otp verify response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return "", "", errors.New(resp.Status)
		}

		// We don't care if this fails
		data := VerifyOtpResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		log.Debug().Msgf("Verify OTP response: %+v", data)
		return data.Data.AccessToken, data.Data.RefreshToken, nil
	}
}
