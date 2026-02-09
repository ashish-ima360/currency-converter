package dto

type ExchangeRateRequest struct {
	FromCurrencyID int     `json:"from_currency_id" binding:"required"`
	ToCurrencyID   int     `json:"to_currency_id" binding:"required"`
	Rate           float64 `json:"rate" binding:"required"`
}

type ExchangeRateResponse struct {
	ID             int     `json:"id"`
	FromCurrencyID int     `json:"from_currency_id"`
	ToCurrencyID   int     `json:"to_currency_id"`
	Rate           float64 `json:"rate"`
	IsActive       bool    `json:"is_active"`
	Deleted        bool    `json:"deleted"`
	DeletedAt      string  `json:"deleted_at"`
	UpdatedAt      string  `json:"updated_at"`
	CreatedAt      string  `json:"created_at"`
}

type ExchangeRateListResponse struct {
	ExchangeRates []ExchangeRateResponse `json:"exchange_rates"`
}