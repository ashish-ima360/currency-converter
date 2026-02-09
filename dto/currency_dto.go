package dto

type CurrencyRequest struct {
	Code   string `json:"code" binding:"required,len=3"`
	Name   string `json:"name" binding:"required"`
	Symbol string `json:"symbol" binding:"required"`
}

type CurrencyResponse struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	IsActive  bool   `json:"is_active"`
	Deleted   bool   `json:"deleted"`
	DeletedAt string `json:"deleted_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

type CurrencyListResponse struct {
	Currencies []CurrencyResponse `json:"currencies"`
}

type CurrencyUpdateRequest struct {
	Name     *string `json:"name"`
	Symbol   *string `json:"symbol"`
	IsActive *bool   `json:"is_active"`
}
