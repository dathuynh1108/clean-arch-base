package repo

import (
	"context"

	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/transaction"
	"gorm.io/gorm"
)

type RepoRouter[Repo any] interface {
	FromMaster(ctx context.Context) Repo
	FromReplica(ctx context.Context) Repo
	Transaction(ctx context.Context, execute func(ctx context.Context, db *gorm.DB) error) error
}

func NewRepoRouter[Repo any](
	dbAlias dbpool.DBAlias,
	dbPool dbpool.DBPool[*gorm.DB],
	repoCreator func(db *gorm.DB) Repo,
) RepoRouter[Repo] {
	return &repoRouter[Repo]{
		dbAlias:     dbAlias,
		dbPool:      dbPool,
		repoCreator: repoCreator,
	}
}

type repoRouter[Repo any] struct {
	dbAlias     dbpool.DBAlias
	dbPool      dbpool.DBPool[*gorm.DB]
	repoCreator func(db *gorm.DB) Repo
}

func (r *repoRouter[Repo]) FromMaster(ctx context.Context) Repo {
	return r.repoCreator(r.dbPool.GetMaster(r.dbAlias).WithContext(ctx))
}

func (r *repoRouter[Repo]) FromReplica(ctx context.Context) Repo {
	return r.repoCreator(r.dbPool.GetMaster(r.dbAlias).WithContext(ctx))
}

func (r *repoRouter[Repo]) Transaction(ctx context.Context, execute func(ctx context.Context, db *gorm.DB) error) error {
	// Check current ctx have transaction or not
	creator := func() *gorm.DB {
		return r.dbPool.GetDB(dbpool.AliasMaster).WithContext(ctx)
	}
	return transaction.Transaction(ctx, execute, creator, nil)
}
