package service

import (
	"context"
	"currency-converter/dto"
	"currency-converter/utils"
	"net/http"
)

type conversionService struct {
	currencyRepo     CurrencyRepository
	exchangeRateRepo ExchangeRateRepository
}

func NewConversionService(currencyRepo CurrencyRepository, exchangeRateRepo ExchangeRateRepository) *conversionService {
	return &conversionService{
		currencyRepo:     currencyRepo,
		exchangeRateRepo: exchangeRateRepo,
	}
}

func (s *conversionService) ConvertCurrency(ctx context.Context, cmd dto.ConversionCmd) (dto.ConversionResult, *utils.AppError) {

	// validate that both currencies exist and are active
	fromCurrency, err := s.currencyRepo.GetByCode(ctx, cmd.From)
	if err != nil {
		return dto.ConversionResult{}, utils.New(http.StatusNotFound, "from currency not found")
	}
	toCurrency, err := s.currencyRepo.GetByCode(ctx, cmd.To)
	if err != nil {
		return dto.ConversionResult{}, utils.New(http.StatusNotFound, "to currency not found")
	}

	if !fromCurrency.IsActive {
		return dto.ConversionResult{}, utils.New(http.StatusBadRequest, "from currency is inactive")
	}

	if !toCurrency.IsActive {
		return dto.ConversionResult{}, utils.New(http.StatusBadRequest, "to currency is inactive")
	}

	// fetch the exchange rate from exchange_rates table
	exchangeRate, err := s.exchangeRateRepo.GetExchangeRateBetweenCurrencies(ctx, fromCurrency.ID, toCurrency.ID)
	if err != nil {
		return dto.ConversionResult{}, utils.New(http.StatusNotFound, "exchange rate not found or inactive")
	}
	
	// convertedAmount = amount * rate
	convertedAmount := cmd.Amount * exchangeRate.Rate

	// return convertedAmount and rate
	return dto.ConversionResult{
		ConvertedAmount: convertedAmount,
		Rate:            exchangeRate.Rate,
	}, nil
}
