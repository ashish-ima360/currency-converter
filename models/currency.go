package models

import (
	"time"
)

type Currency struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Code      string    `gorm:"column:code;not null;size:3;uppercase"`
	Name      string    `gorm:"column:name;not null"`
	Symbol    string    `gorm:"column:symbol;not null"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	Deleted   bool      `gorm:"column:deleted;default:false; not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	DeletedAt time.Time `gorm:"column:deleted_at;autoUpdateTime:false"`
}

// Note: The partial unique index needs to be created via a migration
// Example migration code:
// db.Exec("CREATE UNIQUE INDEX idx_unique_code_not_deleted ON currencies (code) WHERE deleted = false")
