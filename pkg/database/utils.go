package database

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
)

const (
	// Reference: https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	MySQLErrorCodeDuplicateEntry = 1062

	// Reference: https://www.postgresql.org/docs/13/errcodes-appendix.html
	PgSQLErrorCodeDuplicateEntry = "23505"
)

func IsErrorDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}

	err = comerr.UnwrapFirst(err)
	switch errT := err.(type) {
	case *mysql.MySQLError:
		return errT.Number == MySQLErrorCodeDuplicateEntry
	case *pgconn.PgError:
		return errT.Code == PgSQLErrorCodeDuplicateEntry
	default:
		return false
	}
}
