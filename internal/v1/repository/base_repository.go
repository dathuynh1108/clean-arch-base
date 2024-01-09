package repository

import "gorm.io/gorm"

type baseRepository struct {
	*gorm.DB
}
