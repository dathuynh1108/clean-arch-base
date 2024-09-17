package models

import (
	"time"

	"gorm.io/gorm/schema"
)

const (
	DemoTableName = "demo"
)

type Demo struct {
	ID        uint32    `gorm:"id;primaryKey"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (Demo) TableName(namer schema.Namer) string {
	return namer.TableName(DemoTableName)
}
