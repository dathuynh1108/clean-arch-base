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

func ProvideHealthRepo() repo.RepoRouter[HealthRepo] {
	wire.Build(
		ProvideDBAliasDefault,
		database.GetDBPool,
		HealthSet,
	)
	return &healthRepoRouter{}
}
