package sqlquery

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
)

func JoinComma(values ...string) string {
	return strings.Join(values, ", ")
}

func In(colName string, values any) clause.Expression {
	reflectVal := reflect.ValueOf(values)
	switch reflectVal.Kind() {
	case reflect.Array, reflect.Slice:
		break
	default:
		logger.GetLogger().
			WithField("values", values).
			Error("db query IN expression with non-list value")
		return gorm.Expr("0")
	}
	var (
		valueList = make([]any, reflectVal.Len())
	)
	for i := 0; i < reflectVal.Len(); i++ {
		valueList[i] = reflectVal.Index(i).Interface()
	}
	return clause.IN{
		Column: colName,
		Values: valueList,
	}
}

func InWithTableName(tableName string, colName string, values any) clause.Expression {
	reflectVal := reflect.ValueOf(values)
	switch reflectVal.Kind() {
	case reflect.Array, reflect.Slice:
		break
	default:
		logger.GetLogger().
			WithField("values", values).
			Error("db query IN expression with non-list value")
		return gorm.Expr("0")
	}
	var (
		valueList = make([]any, reflectVal.Len())
	)
	for i := 0; i < reflectVal.Len(); i++ {
		valueList[i] = reflectVal.Index(i).Interface()
	}
	return clause.IN{
		Column: gorm.Expr("?.?", clause.Table{Name: tableName}, clause.Column{Name: colName}),
		Values: valueList,
	}
}

func NotIn(colName string, values any) clause.Expression {
	return clause.Not(In(colName, values))
}

func Like(colName any, value any) clause.Expression {
	return clause.Like{
		Column: colName,
		Value:  value,
	}
}

func Gt(colName any, value any) clause.Expression {
	return clause.Gt{
		Column: colName,
		Value:  value,
	}
}

func Gte(colName any, value any) clause.Expression {
	return clause.Gte{
		Column: colName,
		Value:  value,
	}
}

func Equal(colName any, value any) clause.Expression {
	return clause.Eq{
		Column: colName,
		Value:  value,
	}
}

func NotEqual(colName any, value any) clause.Expression {
	return clause.Neq{
		Column: colName,
		Value:  value,
	}
}

func Lte(colName any, value any) clause.Expression {
	return clause.Lte{
		Column: colName,
		Value:  value,
	}
}

func Lt(colName any, value any) clause.Expression {
	return clause.Lt{
		Column: colName,
		Value:  value,
	}
}

func Between(colName string, fromValue any, toValue any) clause.Expr {
	return gorm.Expr("? BETWEEN ? AND ?", clause.Column{Name: colName}, fromValue, toValue)
}

func BetweenTime(colName string, fromTime, toTime int64) clause.Expr {
	return Between(colName, fromTime, toTime-1)
}

func SearchLike(colName string, value string) clause.Expr {
	return gorm.Expr("? LIKE ?", clause.Column{Name: colName}, "%"+value+"%")
}

func OrderAsc(colName string) clause.OrderByColumn {
	return clause.OrderByColumn{
		Column: clause.Column{Name: colName},
		Desc:   false,
	}
}

func OrderAscWithTableName(tableName, colName string) clause.OrderByColumn {
	return clause.OrderByColumn{
		Column: clause.Column{Table: tableName, Name: colName},
		Desc:   false,
	}
}

func OrderAscF(colName string) clause.OrderByColumn {
	order := OrderAsc(colName)
	order.Reorder = true
	return order
}

func OrderAscEx(colName string) ExOrderBy {
	return ExOrderBy{
		OrderBy: clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: Col(colName),
				},
			},
		},
	}
}

func OrderDesc(colName string) clause.OrderByColumn {
	return clause.OrderByColumn{
		Column: clause.Column{Name: colName},
		Desc:   true,
	}
}

func OrderDescWithTableName(tableName, colName string) clause.OrderByColumn {
	return clause.OrderByColumn{
		Column: clause.Column{Table: tableName, Name: colName},
		Desc:   true,
	}
}

func OrderDescF(colName string) clause.OrderByColumn {
	order := OrderDesc(colName)
	order.Reorder = true
	return order
}

func OrderDescEx(colName string) ExOrderBy {
	return ExOrderBy{
		OrderBy: clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: Col(colName),
					Desc:   true,
				},
			},
		},
	}
}

func OrderRecentPK() clause.OrderByColumn {
	return clause.OrderByColumn{
		Column:  clause.PrimaryColumn,
		Desc:    true,
		Reorder: true,
	}
}

func OrderDescPK() clause.OrderByColumn {
	return clause.OrderByColumn{
		Column:  clause.PrimaryColumn,
		Desc:    true,
		Reorder: true,
	}
}

func DbOrderEx(db *gorm.DB, tableName, colName string, desc bool) *gorm.DB {
	var order clause.OrderByColumn
	if desc {
		order = OrderDescWithTableName(tableName, colName)
	} else {
		order = OrderAscWithTableName(tableName, colName)
	}
	return db.Clauses(ExOrderBy{
		OrderBy: clause.OrderBy{
			Columns: []clause.OrderByColumn{order},
		},
	})
}

func OrderExpression(expr clause.Expression, desc bool) ExOrderBy {
	return ExOrderBy{
		OrderBy: clause.OrderBy{
			Expression: expr,
		},
		IsExpressionDesc: desc,
	}
}

func MaxAlias(colName string, alias string) clause.Expr {
	return gorm.Expr("MAX(?) AS ?", Col(colName), Col(alias))
}

func Max(colName string) clause.Expr {
	return MaxAlias(colName, "max_"+colName)
}

func MinAlias(colName string, alias string) clause.Expr {
	return gorm.Expr("MIN(?) AS ?", clause.Column{Name: colName}, clause.Column{Name: alias})
}

func Min(colName string) clause.Expr {
	return MinAlias(colName, "min_"+colName)
}

func SumAlias(colName string, alias string) clause.Expr {
	return gorm.Expr("SUM(?) AS ?", Col(colName), Col(alias))
}

func Sum(colName string) clause.Expr {
	return SumAlias(colName, "sum_"+colName)
}

func SumGtAlias(colName string, gtValue any, alias string) clause.Expr {
	return gorm.Expr("SUM(CASE WHEN ? > ? THEN ? ELSE 0 END) AS ?",
		Col(colName), gtValue, Col(colName), Col(alias))
}

func SumGt(colName string, gtValue any) clause.Expr {
	return SumGtAlias(colName, gtValue, "sum_gt_"+colName)
}

func SumLtAlias(colName string, ltValue any, alias string) clause.Expr {
	return gorm.Expr("SUM(CASE WHEN ? < ? THEN ? ELSE 0 END) AS ?",
		Col(colName), ltValue, Col(colName), Col(alias))
}

func SumLt(colName string, ltValue any) clause.Expr {
	return SumLtAlias(colName, ltValue, "sum_lt_"+colName)
}

func SumIfAlias(colName string, expr clause.Expression, alias string) clause.Expr {
	return gorm.Expr("SUM(IF(?, ?, 0)) as ?", expr, Col(colName), Col(alias))
}

func SumIf(colName string, expr clause.Expression) clause.Expr {
	return SumIfAlias(colName, expr, "sum_eq_"+colName)
}

func CountAlias(colName string, alias string) clause.Expr {
	return gorm.Expr("COUNT(?) AS ?", Col(colName), Col(alias))
}

func Count(colName string) clause.Expr {
	return CountAlias(colName, "count_"+colName)
}

func CountDistinct(colName string) clause.Expr {
	return CountDistinctAlias(colName, "count_distinct_"+colName)
}

func CountDistinctAlias(colName string, alias string) clause.Expr {
	return gorm.Expr("COUNT(DISTINCT ?) AS ?", Col(colName), Col(alias))
}

func CountIfAlias(expr clause.Expression, alias string) clause.Expr {
	return gorm.Expr("COUNT(IF(?, id, null)) as ?", expr, Col(alias))
}

func CountAll() clause.Expr {
	return gorm.Expr("COUNT(*)")
}

func CountAllAlias(alias string) clause.Expr {
	return gorm.Expr("COUNT(*) AS ?", Col(alias))
}

func Table(tableName string) clause.Table {
	return clause.Table{Name: tableName}
}

func Col(colName string) clause.Column {
	return clause.Column{Name: colName}
}

func ColWithTableName(tableName, colName string) clause.Column {
	return clause.Column{
		Table: tableName,
		Name:  colName,
	}
}

func ColWithTableNameAlias(tableName, colName, alias string) clause.Column {
	return clause.Column{
		Table: tableName,
		Name:  colName,
		Alias: alias,
	}
}

func ColWithTableNameByValue(tableName, colName string, value any) clause.Expr {
	return gorm.Expr("?.? = ?", Table(tableName), Col(colName), value)
}

func ColIncr(colName string) clause.Expr {
	return gorm.Expr("?+1", Col(colName))
}

func ColDecr(colName string) clause.Expr {
	return gorm.Expr("?-1", Col(colName))
}

func ColAdd(colName string, value any) clause.Expr {
	return gorm.Expr("?+?", Col(colName), value)
}

func ColSub(colName string, value any) clause.Expr {
	return gorm.Expr("?-?", Col(colName), value)
}

func DbSelectForUpdate(db *gorm.DB) *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}

func DbNot(db *gorm.DB, query any, args ...any) (tx *gorm.DB) {
	tx = db.Clauses() // Like `db.getInstance()`
	conds := tx.Statement.BuildCondition(query, args...)
	if len(conds) == 0 {
		return db
	}
	if len(conds) == 1 {
		tx.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{clause.Not(conds[0])},
		})
	} else {
		notClauses := make([]clause.Expression, len(conds))
		for i, cond := range conds {
			notClauses[i] = clause.Not(cond)
		}
		tx.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{clause.Or(notClauses...)},
		})
	}
	return tx
}

func InSubQuery(colName string, query *gorm.DB) clause.Expression {
	if query == nil {
		return gorm.Expr("0")
	}
	return gorm.Expr("? IN (?)", Col(colName), query)
}

func NotInSubQuery(colName string, query *gorm.DB) clause.Expression {
	if query == nil {
		return gorm.Expr("0")
	}
	return gorm.Expr("? NOT IN (?)", Col(colName), query)
}

// unixTime can be a time value or a column
func FromUnixTime(unixTime interface{}, format string) clause.Expression {
	if format == "" {
		return gorm.Expr("from_unixtime(?)", unixTime)
	}
	return gorm.Expr("from_unixtime(?, ?)", unixTime, format)
}

func UnixTimeToDateAlias(colName string, format, alias string) clause.Expr {
	return gorm.Expr("DATE(?) as ?", FromUnixTime(Col(colName), format), Col(alias))
}

func OrderAscPositiveFirst(colName string) ExOrderBy {
	return ExOrderBy{
		OrderBy: clause.OrderBy{
			Expression: gorm.Expr("IF(? > 0, 0, 1), ABS(?)", Col(colName), Col(colName)),
		},
		IsExpressionDesc: false,
	}
}
