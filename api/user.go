package api

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/rs/zerolog/log"
)

type UserResponse struct {
	Meta struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
		URI       string    `json:"uri"`
	} `json:"meta"`
	Data struct {
		User struct {
			UUID           string      `json:"uuid"`
			FirstName      string      `json:"first_name"`
			MiddleName     interface{} `json:"middle_name"`
			LastName       string      `json:"last_name"`
			Email          string      `json:"email"`
			EmailVerified  bool        `json:"email_verified"`
			Phone          string      `json:"phone"`
			PhoneVerified  bool        `json:"phone_verified"`
			GoogleLinked   bool        `json:"google_linked"`
			AppleLinked    bool        `json:"apple_linked"`
			Role           string      `json:"role"`
			IsInternalUser bool        `json:"is_internal_user"`
			BetaAccess     bool        `json:"beta_access"`
			WebBetaAccess  bool        `json:"web_beta_access"`
			CcEnabled      bool        `json:"cc_enabled"`
			Timezone       string      `json:"timezone"`
			CreatedAt      time.Time   `json:"created_at"`
			UpdatedAt      time.Time   `json:"updated_at"`
		} `json:"user"`
		Settings struct {
			HasActiveConsent bool `json:"has_active_consent"`
			Widgets          struct {
				BankBalance struct {
					ShowTotal       bool        `json:"show_total"`
					DoubleTapToHide bool        `json:"double_tap_to_hide"`
					WidgetEnabled   bool        `json:"widget_enabled"`
					Order           interface{} `json:"order"`
					BankAccounts    []struct {
						AccountID           string `json:"account_id"`
						AccountName         string `json:"account_name"`
						MaskedAccountNumber string `json:"masked_account_number"`
						ShowInWidget        bool   `json:"show_in_widget"`
						Order               int    `json:"order"`
					} `json:"bank_accounts"`
				} `json:"bank_balance"`
				CashFlow struct {
					ShowIncoming    bool        `json:"show_incoming"`
					ShowOutgoing    bool        `json:"show_outgoing"`
					ShowInvested    bool        `json:"show_invested"`
					ShowLeftToSpend bool        `json:"show_left_to_spend"`
					DoubleTapToHide bool        `json:"double_tap_to_hide"`
					FieldOrder      []string    `json:"field_order"`
					WidgetEnabled   bool        `json:"widget_enabled"`
					Order           interface{} `json:"order"`
					BankAccounts    []struct {
						AccountID           string `json:"account_id"`
						AccountName         string `json:"account_name"`
						MaskedAccountNumber string `json:"masked_account_number"`
						ShowInWidget        bool   `json:"show_in_widget"`
					} `json:"bank_accounts"`
				} `json:"cash_flow"`
				BankAccount struct {
					WidgetEnabled bool        `json:"widget_enabled"`
					Order         interface{} `json:"order"`
					BankAccounts  []struct {
						AccountID           string `json:"account_id"`
						AccountName         string `json:"account_name"`
						MaskedAccountNumber string `json:"masked_account_number"`
						ShowInWidget        bool   `json:"show_in_widget"`
						Order               int    `json:"order"`
					} `json:"bank_accounts"`
				} `json:"bank_account"`
				OtherBankAccount struct {
					WidgetEnabled bool        `json:"widget_enabled"`
					Order         interface{} `json:"order"`
					BankAccounts  interface{} `json:"bank_accounts"`
				} `json:"other_bank_account"`
				SpendingSummary struct {
					WidgetEnabled bool        `json:"widget_enabled"`
					TagsType      string      `json:"tags_type"`
					CustomTags    interface{} `json:"custom_tags"`
					DataDuration  string      `json:"data_duration"`
					DataFormat    string      `json:"data_format"`
					BankAccounts  []struct {
						AccountID           string `json:"account_id"`
						AccountName         string `json:"account_name"`
						MaskedAccountNumber string `json:"masked_account_number"`
						ShowInWidget        bool   `json:"show_in_widget"`
					} `json:"bank_accounts"`
				} `json:"spending_summary"`
			} `json:"widgets"`
		} `json:"settings"`
		Onboarding struct {
			Status                   bool `json:"status"`
			WelcomeScreen            bool `json:"welcome_screen"`
			RefreshExplainer         bool `json:"refresh_explainer"`
			BankAccountLinked        bool `json:"bank_account_linked"`
			FetchingDataScreen       bool `json:"fetching_data_screen"`
			NotifPermissionExplainer bool `json:"notif_permission_explainer"`
		} `json:"onboarding"`
		Route string `json:"route"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

func User() (string, error) {

	RefreshOrFail()

	req, _ := APIRequest("GET", Url("/v2/users/me"), nil)
	resp, err := Client.Do(req)

	if err != nil {
		return "", err
	} else {

		log.Debug().Msgf("User response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return "", errors.New(resp.Status)
		}

		data := UserResponse{}
		json.NewDecoder(resp.Body).Decode(&data)
		log.Debug().Msgf("Fetch user response: %+v", data)

		return data.Data.User.UUID, nil
	}
}
