package models

import "time"

type ExchangeRate struct {
	ID             int       `gorm:"column:id;primaryKey;autoIncrement"`
	FromCurrencyID int       `gorm:"column:from_currency_id;not null;reference:currencies(id)"`
	ToCurrencyID   int       `gorm:"column:to_currency_id;not null;reference:currencies(id)"`
	Rate           float64   `gorm:"column:rate;not null"`
	IsActive       bool      `gorm:"column:is_active;default:true"`
	Deleted        bool      `gorm:"column:deleted;default:false;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	DeletedAt      time.Time `gorm:"column:deleted_at;autoUpdateTime:false"`
}
