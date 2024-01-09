package repository

import "gorm.io/gorm"

type HealthRepository struct {
	baseRepository
}

func NewHealthRepository(db *gorm.DB) *HealthRepository {
	return &HealthRepository{
		baseRepository: baseRepository{db},
	}
}
