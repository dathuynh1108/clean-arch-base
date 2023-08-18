package repository

import "gorm.io/gorm"

type HealthRepository struct {
	db *gorm.DB
}

func NewHealthRepository(db *gorm.DB) *HealthRepository {
	return &HealthRepository{db}
}
