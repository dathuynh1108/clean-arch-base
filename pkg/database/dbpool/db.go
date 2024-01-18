package dbpool

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{db}
}

func (db *DB) WithContext(ctx context.Context) *DB {
	return &DB{db.DB.WithContext(ctx)}
}

func (db *DB) WithLogger(logger logger.Interface) *DB {
	return &DB{
		db.Session(&gorm.Session{Logger: logger}),
	}
}
