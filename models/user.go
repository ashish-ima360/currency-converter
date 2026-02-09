package models

import "time"

type User struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	Email        string    `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
}
