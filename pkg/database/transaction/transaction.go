package transaction

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

func Transaction(
	ctx context.Context,
	execute func(context.Context, *gorm.DB) error,
	creator func() *gorm.DB,
	options *sql.TxOptions,
) (err error) {
	// Check current ctx have transaction or not
	// If not, create new transaction
	// If yes, use current transaction
	var tx *gorm.DB
	var executeCtx context.Context
	if ctx.Value(ContextKeyTransaction) == nil {
		tx = creator()
		executeCtx = context.WithValue(ctx, ContextKeyTransaction, tx)
	} else {
		tx = ctx.Value(ContextKeyTransaction).(*gorm.DB)
		executeCtx = ctx
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = execute(executeCtx, tx)
	if err == nil {
		err = tx.Commit().Error
	}
	return
}
