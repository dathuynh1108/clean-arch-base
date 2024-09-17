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
		wire.Bind(new(repo.RepoRouter[HealthRepo]), new(*healthRepoRouter)),
	)
)

type HealthRepo interface {
	GetDemo() (any, error)
}

type healthRepo struct {
	baseRepository
}

func NewHealthRepo(db *gorm.DB) *healthRepo {
	return &healthRepo{
		baseRepository: baseRepository{db},
	}
}

type healthRepoRouter struct {
	repo.RepoRouter[HealthRepo]
}

func NewHealthRepoRouter(alias dbpool.DBAlias, dbPool dbpool.DBPool) *healthRepoRouter {
	return &healthRepoRouter{
		repo.NewRepoRouter(
			alias,
			dbPool,
			func(db *gorm.DB) HealthRepo {
				return NewHealthRepo(db)
			},
		),
	}
}

func (r *healthRepo) GetDemo() (res any, err error) {
	err = r.DB.Take(&res).Error
	return
}
