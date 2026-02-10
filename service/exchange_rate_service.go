package service

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/utils"
	"encoding/json"
	"net/http"
)

type ExchangeRateRepository interface {
	Create(ctx context.Context, exchangeRate *models.ExchangeRate) (*models.ExchangeRate, error)
	GetByID(ctx context.Context, id int) (*models.ExchangeRate, error)
	GetAll(ctx context.Context) ([]models.ExchangeRate, error)
	Update(ctx context.Context, exchangeRate *models.ExchangeRate) error
	Delete(ctx context.Context, id int) error
	GetExchangeRateBetweenCurrencies(ctx context.Context, fromCurrencyID int, toCurrencyID int) (models.ExchangeRate, error)
	CreateOrUpdate(ctx context.Context, fromCurrencyID int, toCurrencyID int, rate float64) error
}

type exchangeRateService struct {
	repo            ExchangeRateRepository
	currencyRepo    CurrencyRepository
	httpClient      *http.Client
	exchangeRateAPI string
}

func NewExchangeRateService(
	repo ExchangeRateRepository,
	currencyRepo CurrencyRepository,
	httpClient *http.Client,
	exchangeRateAPI string,
	) *exchangeRateService {
	return &exchangeRateService{
		repo:            repo,
		currencyRepo:    currencyRepo,
		httpClient:      httpClient,
		exchangeRateAPI: exchangeRateAPI,
	}
}

func (s *exchangeRateService) CreateExchangeRate(ctx context.Context, req dto.ExchangeRateRequest) (*models.ExchangeRate, *utils.AppError) {
	exchangeRate := &models.ExchangeRate{
		FromCurrencyID: req.FromCurrencyID,
		ToCurrencyID:   req.ToCurrencyID,
		Rate:           req.Rate,
	}

	createdExchangeRate, err := s.repo.Create(ctx, exchangeRate)
	if err != nil {
		return nil, utils.New(http.StatusInternalServerError, "error in creating exchange rate")
	}

	return createdExchangeRate, nil
}

func (s *exchangeRateService) GetExchangeRateByID(ctx context.Context, id int) (*models.ExchangeRate, *utils.AppError) {
	exchangeRate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.New(http.StatusNotFound, "exchange rate not found")
	}
	return exchangeRate, nil
}

func (s *exchangeRateService) GetAllExchangeRates(ctx context.Context) ([]models.ExchangeRate, *utils.AppError) {
	exchangeRates, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, utils.New(http.StatusInternalServerError, "error in fetching all exchange rates")
	}
	return exchangeRates, nil
}

func (s *exchangeRateService) UpdateExchangeRate(ctx context.Context, id int, req dto.ExchangeRateUpdateRequest) (*models.ExchangeRate, *utils.AppError) {
	exchangeRate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.New(http.StatusNotFound, "exchange rate not found")
	}
	if req.Rate != nil {
		exchangeRate.Rate = *req.Rate
	}
	if req.IsActive != nil {
		exchangeRate.IsActive = *req.IsActive
	}

	err = s.repo.Update(ctx, exchangeRate)
	if err != nil {
		return nil, utils.New(http.StatusInternalServerError, "error in updating exchange rate")
	}

	return exchangeRate, nil
}

func (s *exchangeRateService) DeleteExchangeRate(ctx context.Context, id int) *utils.AppError {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return utils.New(http.StatusNotFound, "exchange rate not found")
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return utils.New(http.StatusInternalServerError, "error in deleting exchange rate")
	}
	return nil
}

func (s *exchangeRateService) SyncExchangeRates(ctx context.Context, code string) *utils.AppError {
	// validation done in controller

	// build a get request to fetch exchange rates for the given code from the external API
	url := s.exchangeRateAPI + "/" + code

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return utils.New(http.StatusInternalServerError, "error in creating http request")
	}
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return utils.New(http.StatusInternalServerError, "error in making http request to exchange rate API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return utils.New(http.StatusInternalServerError, "received non-200 response from exchange rate API")
	}

	// parse the response and update the exchange rates in the database using the repo
	var apiResponse dto.ExchangeRateExternalResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return utils.New(http.StatusInternalServerError, "error in parsing response from exchange rate API")
	}
	if apiResponse.Result != "success" {
		return utils.New(http.StatusInternalServerError, "exchange rate API returned unsuccessful result")
	}
	if apiResponse.BaseCode != code {
		return utils.New(http.StatusInternalServerError, "exchange rate API returned data for unexpected base code")
	}

	fromCurrency, err := s.currencyRepo.GetByCode(ctx, apiResponse.BaseCode)
	if err != nil {
		return utils.New(http.StatusInternalServerError, "error in fetching from currency ID")
	}

	toCurrencyCodes := []string{"USD", "INR", "EUR", "CAD", "JPY"}

	for _, toCurrencyCode := range toCurrencyCodes {
		rate, ok := apiResponse.ConversionRates[toCurrencyCode]
		if !ok || toCurrencyCode == code {
			continue // skip if the API response does not contain conversion rate for this currency code
		}
		toCurrency, err := s.currencyRepo.GetByCode(ctx, toCurrencyCode)
		if err != nil {
			return utils.New(http.StatusInternalServerError, "error in fetching to currency ID")
		}

		// update the exchange rate in the database
		err = s.repo.CreateOrUpdate(ctx, fromCurrency.ID, toCurrency.ID, rate)
		if err != nil {
			return utils.New(http.StatusInternalServerError, "error in updating exchange rate in database")
		}
	}

	return nil
}
