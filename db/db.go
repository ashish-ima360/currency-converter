package db

import (
	"currency-converter/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dbUrl string) (*gorm.DB, error) {

	// log.Printf("DNS is : %v", dbUrl)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {

	if err := db.AutoMigrate(&models.User{}, &models.Currency{}, &models.ExchangeRate{}); err != nil {
		return err
	}
	
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_code_not_deleted 
		ON currencies (code) 
		WHERE deleted = false
	`).Error; err != nil {
		return err
	}

	return db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_uique_exchange_rate_not_deleted
		ON exchange_rates (from_currency_id, to_currency_id)
		WHERE deleted = false
	`).Error
}
