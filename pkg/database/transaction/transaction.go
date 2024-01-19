package transaction

import (
	"context"
	"database/sql"

	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"gorm.io/gorm"
)

func Transaction(
	ctx context.Context,
	alias dbpool.DBAlias,
	execute func(context.Context, *gorm.DB) error,
	options *sql.TxOptions,
) (err error) {
	// Check current ctx have transaction or not
	// If not, create new transaction
	// If yes, use current transaction, gorm Transaction() handle nested transaction
	var (
		tx         *gorm.DB
		executeCtx context.Context
	)
	if ctx.Value(ContextKeyTransaction) == nil {
		tx = database.ProvideDBPool().GetMaster(alias).WithContext(ctx).DB
		executeCtx = context.WithValue(ctx, ContextKeyTransaction, tx)
	} else {
		tx = ctx.Value(ContextKeyTransaction).(*gorm.DB)
		executeCtx = ctx
	}

	// Pass executeCtx to transaction execution
	return tx.Transaction(func(tx *gorm.DB) error {
		return execute(executeCtx, tx)
	}, options)
}
