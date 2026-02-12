package repository

import (
	"context"
	"currency-converter/dto"
	"currency-converter/models"
	"currency-converter/utils"
	"time"

	"gorm.io/gorm"
)

type currencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) *currencyRepository {
	return &currencyRepository{
		db: db,
	}
}

func (r *currencyRepository) Create(ctx context.Context, currency *models.Currency) (*models.Currency, error) {

	err := r.db.WithContext(ctx).Create(currency).Error
	if err != nil {
		return nil, err
	}
	return currency, nil
}

func (r *currencyRepository) GetByID(ctx context.Context, id int) (*models.Currency, error) {
	var currency models.Currency

	err := r.db.WithContext(ctx).Where("id = ? AND deleted = ?", id, false).First(&currency).Error
	if err != nil {
		return nil, err
	}
	return &currency, nil
}

func (r *currencyRepository) GetAll(ctx context.Context) ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.WithContext(ctx).Where("deleted = ?", false).Find(&currencies).Error
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

func (r *currencyRepository) Update(ctx context.Context, id int, input dto.CurrencyUpdateRequest) error {

	tx := r.db.WithContext(ctx).
		Model(&models.Currency{}).
		Where("id = ?", id).
		Updates(input).Update("updated_at", time.Now())

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return utils.ErrCodeNotFound
	}

	return nil
}

func (r *currencyRepository) Delete(ctx context.Context, id int) error {

	return r.db.WithContext(ctx).
		Model(&models.Currency{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted":    true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *currencyRepository) GetByCode(ctx context.Context, code string) (models.Currency, error) {
	var currency models.Currency
	err := r.db.WithContext(ctx).Where("code = ? AND deleted = ?", code, false).First(&currency).Error
	if err != nil {
		return models.Currency{}, err
	}
	return currency, nil
}
