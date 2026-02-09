package repository

import (
	"context"
	"currency-converter/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type exchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) *exchangeRateRepository {
	return &exchangeRateRepository{
		db: db,
	}
}

func (r *exchangeRateRepository) Create(ctx context.Context, exchangeRate *models.ExchangeRate) (*models.ExchangeRate, error) {

	err := r.db.WithContext(ctx).Create(&exchangeRate).Error
	if err != nil {
		return nil, err
	}
	return exchangeRate, nil
}

func (r *exchangeRateRepository) GetByID(ctx context.Context, id int) (*models.ExchangeRate, error) {
	var exchangeRate models.ExchangeRate

	err := r.db.WithContext(ctx).Where("id = ? AND deleted = ?", id, false).First(&exchangeRate).Error
	if err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func (r *exchangeRateRepository) GetAll(ctx context.Context) ([]models.ExchangeRate, error) {
	var exchangeRates []models.ExchangeRate

	err := r.db.WithContext(ctx).Where("deleted = ?", false).Find(&exchangeRates).Error
	if err != nil {
		return nil, err
	}
	return exchangeRates, nil
}

func (r *exchangeRateRepository) Update(ctx context.Context, id int, rate float64) error {
	result := r.db.WithContext(ctx).Model(&models.ExchangeRate{}).Where("id = ?", id).Update("rate", rate)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("exchange rate not found")
	}
	return nil
}

func (r *exchangeRateRepository) Delete(ctx context.Context, id int) error {

	return r.db.WithContext(ctx).
		Model(&models.ExchangeRate{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted": true, 
			"deleted_at": time.Now(),
		}).Error
}

func (r *exchangeRateRepository) GetExchangeRateBetweenCurrencies(ctx context.Context, fromCurrencyID int, toCurrencyID int) (models.ExchangeRate, error) {
	var exchangeRate models.ExchangeRate
	err := r.db.WithContext(ctx).Where("from_currency_id = ? AND to_currency_id = ? AND is_active = ? AND deleted = ?", fromCurrencyID, toCurrencyID, true, false).
		First(&exchangeRate).Error
	if err != nil {
		return models.ExchangeRate{}, err
	}
	return exchangeRate, nil
}