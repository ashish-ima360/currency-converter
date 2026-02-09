package service

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/utils"
	"net/http"
)

type ExchangeRateRepository interface {
	Create(ctx context.Context, exchangeRate *models.ExchangeRate) (*models.ExchangeRate, error)
	GetByID(ctx context.Context, id int) (*models.ExchangeRate, error)
	GetAll(ctx context.Context) ([]models.ExchangeRate, error)
	Update(ctx context.Context, id int, rate float64) error
	Delete(ctx context.Context, id int) error
	GetExchangeRateBetweenCurrencies(ctx context.Context, fromCurrencyID int, toCurrencyID int) (models.ExchangeRate, error)
}

type exchangeRateService struct {
	repo ExchangeRateRepository
}

func NewExchangeRateService(repo ExchangeRateRepository) *exchangeRateService {
	return &exchangeRateService{
		repo: repo,
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

func (s *exchangeRateService) UpdateExchangeRate(ctx context.Context, id int, rate float64) (*models.ExchangeRate, *utils.AppError) {
	exchangeRate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.New(http.StatusNotFound, "exchange rate not found")
	}
	err = s.repo.Update(ctx, id, rate)
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