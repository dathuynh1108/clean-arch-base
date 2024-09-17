package sqlquery

import (
	"github.com/dathuynh1108/clean-arch-base/pkg/meta"

	"gorm.io/gorm"
)

const (
	PagingEmptyLimit = 500
)

func DBPaging(db *gorm.DB, paging *meta.Paging) *gorm.DB {
	if paging.Limit > 0 {
		db = db.Limit(paging.Limit)
	} else {
		db = db.Limit(PagingEmptyLimit)
	}
	if paging.Offset > 0 {
		db = db.Offset(paging.Offset)
	}
	return db
}
