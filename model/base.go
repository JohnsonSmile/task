package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int       `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

type GormList []string

func (l GormList) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *GormList) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), l)
}

func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
