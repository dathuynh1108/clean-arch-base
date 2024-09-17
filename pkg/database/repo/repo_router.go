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
	Transaction(ctx context.Context, execute func(ctx context.Context, tx Repo) error) error
}

func NewRepoRouter[Repo any](
	dbAlias dbpool.DBAlias,
	dbPool dbpool.DBPool,
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
	dbPool      dbpool.DBPool
	repoCreator func(db *gorm.DB) Repo
}

func (r *repoRouter[Repo]) FromMaster(ctx context.Context) Repo {
	return r.repoCreator(r.dbPool.GetMaster(r.dbAlias).WithContext(ctx).DB)
}

func (r *repoRouter[Repo]) FromReplica(ctx context.Context) Repo {
	return r.repoCreator(r.dbPool.GetReplica(r.dbAlias).WithContext(ctx).DB)
}

func (r *repoRouter[Repo]) Transaction(ctx context.Context, execute func(ctx context.Context, tx Repo) error) error {
	// Check current ctx have transaction or not
	wrapRepo := func(ctx context.Context, tx *gorm.DB) error {
		txRepo := r.repoCreator(tx)
		return execute(ctx, txRepo)
	}
	return transaction.Transaction(ctx, r.dbAlias, wrapRepo)
}
