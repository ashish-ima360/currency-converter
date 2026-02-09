package dto

type CurrencyConversionResponse struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	Amount          float64 `json:"amount"`
	Rate            float64 `json:"rate"`
	ConvertedAmount float64 `json:"converted_amount"`
}

type ConversionCmd struct {
	From   string
	To     string
	Amount float64
}

type ConversionResult struct {
	ConvertedAmount float64
	Rate            float64
}
