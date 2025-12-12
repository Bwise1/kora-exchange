package fxrates

import "time"

// FastForexAPIResponse represents the response from FastForex API
type FastForexAPIResponse struct {
	Base    string             `json:"base"`
	Results map[string]float64 `json:"results"`
	Updated string             `json:"updated"`
	MS      int                `json:"ms"`
}

// Deprecated: ExchangeRateAPIResponse - keeping for reference
type ExchangeRateAPIResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

// FXRate represents a foreign exchange rate
type FXRate struct {
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Rate         float64   `json:"rate"`
	LastUpdated  time.Time `json:"last_updated"`
}

// FXRatesResponse represents the response sent to clients
type FXRatesResponse struct {
	BaseCurrency string             `json:"base_currency"`
	Rates        map[string]float64 `json:"rates"`
	LastUpdated  time.Time          `json:"last_updated"`
}

// ConversionRequest represents a currency conversion request
type ConversionRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// ConversionResponse represents a currency conversion response
type ConversionResponse struct {
	From        string    `json:"from"`
	To          string    `json:"to"`
	Amount      float64   `json:"amount"`
	Result      float64   `json:"result"`
	Rate        float64   `json:"rate"`
	LastUpdated time.Time `json:"last_updated"`
}
