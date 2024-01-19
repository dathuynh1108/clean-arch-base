//go:build wireinject

package repository

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/repo"
	"github.com/google/wire"
)

func ProvideDBAliasDefault() dbpool.DBAlias {
	return dbpool.DBDefault
}

func ProvideHealthRepo() repo.RepoRouter[HealthRepository] {
	wire.Build(
		ProvideDBAliasDefault,
		database.ProvideDBPool,
		HealthSet,
	)
	return &healthRepoRouter{}
}
