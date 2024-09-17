package transaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/pkg/database"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"

	"gorm.io/gorm"
)

func Transaction(
	ctx context.Context,
	alias dbpool.DBAlias,
	execute func(context.Context, *gorm.DB) error,
	opts ...*sql.TxOptions,
) (err error) {
	// Check current ctx have transaction or not
	// If not, create new transaction
	// If yes, use current transaction, gorm Transaction() handle nested transaction
	var (
		tx  *gorm.DB
		key = buildTxContextKey(alias)
	)
	if txCtx := ctx.Value(key); txCtx != nil {
		tx = txCtx.(*gorm.DB)
	} else {
		tx = database.GetDBPool().GetMaster(alias).WithContext(ctx).DB
	}

	// Pass executeCtx to transaction execution
	return tx.Transaction(func(tx *gorm.DB) error {
		return execute(context.WithValue(ctx, key, tx), tx)
	}, opts...)
}

func buildTxContextKey(alias dbpool.DBAlias) ContextKeyTransactionType {
	return ContextKeyTransactionType(fmt.Sprintf("transaction_%s", alias))
}
