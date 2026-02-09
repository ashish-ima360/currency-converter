package service

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/utils"
	"net/http"
	"strings"
)

type CurrencyRepository interface {
	Create(ctx context.Context, currency *models.Currency) (*models.Currency, error)
	GetByID(ctx context.Context, id int) (*models.Currency, error)
	GetAll(ctx context.Context) ([]models.Currency, error)
	Update(ctx context.Context, currency *models.Currency) error
	Delete(ctx context.Context, id int) error
	GetByCode(ctx context.Context, code string) (models.Currency, error)
}

type currencyService struct {
	currencyRepo CurrencyRepository
}

func NewCurrencyService(currencyRepo CurrencyRepository) *currencyService {
	return &currencyService{
		currencyRepo: currencyRepo,
	}
}

func (s *currencyService) CreateCurrency(ctx context.Context, req dto.CurrencyRequest) (*models.Currency, *utils.AppError) {
	currency := &models.Currency{
		Code: strings.ToUpper(req.Code),
		Name: req.Name,
		Symbol: req.Symbol,
	}
	createdCurrency, err := s.currencyRepo.Create(ctx, currency)
	if err != nil {
		return nil, utils.New(http.StatusInternalServerError, "error in creating currency")
	}
	return createdCurrency, nil
}

func (s *currencyService) GetCurrencyByID(ctx context.Context, id int) (*models.Currency, *utils.AppError) {
	currency, err := s.currencyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.New(http.StatusNotFound, "currency not found")
	}
	return currency, nil
}

func (s *currencyService) GetAllCurrencies(ctx context.Context) ([]models.Currency, *utils.AppError) {
	
	currencies, err := s.currencyRepo.GetAll(ctx)
	if err != nil {
		return nil, utils.New(http.StatusInternalServerError, "error in fetching currencies")
	}
	
	return currencies, nil
}

func (s *currencyService) UpdateCurrency(ctx context.Context, id int, req dto.CurrencyUpdateRequest) (*models.Currency, *utils.AppError) {
	
	currency, err := s.currencyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.New(http.StatusNotFound, "currency not found")
	}

	if req.Name != nil {
		currency.Name = *req.Name
	}
	if req.Symbol != nil {
		currency.Symbol = *req.Symbol
	}
	if req.IsActive != nil {
		currency.IsActive = *req.IsActive
	}

	err = s.currencyRepo.Update(ctx, currency)
	if err != nil {
		return nil, utils.New(http.StatusInternalServerError, "error in updating currency")
	}

	return currency, nil
}

func (s *currencyService) DeleteCurrency(ctx context.Context, id int) *utils.AppError {
	
	err := s.currencyRepo.Delete(ctx, id)
	if err != nil {
		return utils.New(http.StatusInternalServerError, "error in deleting currency")
	}
	return nil
}
