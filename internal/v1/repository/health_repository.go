package repository

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/repo"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	HealthSet = wire.NewSet(
		NewHealthRepoRouter,
		wire.Bind(new(repo.RepoRouter[HealthRepository]), new(*healthRepoRouter)),
	)
)

type HealthRepository interface{}

type healthRepository struct {
	baseRepository
}

func NewHealthRepository(db *gorm.DB) *healthRepository {
	return &healthRepository{
		baseRepository: baseRepository{db},
	}
}

type healthRepoRouter struct {
	repo.RepoRouter[HealthRepository]
}

func NewHealthRepoRouter(alias dbpool.DBAlias, dbPool dbpool.DBPool) *healthRepoRouter {
	return &healthRepoRouter{
		repo.NewRepoRouter(
			alias,
			dbPool,
			func(db *gorm.DB) HealthRepository {
				return NewHealthRepository(db)
			},
		),
	}
}
