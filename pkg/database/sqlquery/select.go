package sqlquery

import (
	"strings"

	"gorm.io/gorm/clause"
)

type SelectBuilder struct {
	sqlList []string
	vars    []any
}

func (b SelectBuilder) AddColumn(name string) SelectBuilder {
	return SelectBuilder{
		sqlList: append(b.sqlList, "?"),
		vars:    append(b.vars, Col(name)),
	}
}

func (b SelectBuilder) AddColumns(names ...string) SelectBuilder {
	var (
		addSqlList = make([]string, len(names))
		addColumns = make([]any, len(names))
	)
	for i, name := range names {
		addSqlList[i] = "?"
		addColumns[i] = Col(name)
	}
	return SelectBuilder{
		sqlList: append(b.sqlList, addSqlList...),
		vars:    append(b.vars, addColumns...),
	}
}

func (b SelectBuilder) AddExpr(expr clause.Expr) SelectBuilder {
	return SelectBuilder{
		sqlList: append(b.sqlList, expr.SQL),
		vars:    append(b.vars, expr.Vars...),
	}
}

func (b SelectBuilder) AddExprs(exprs ...clause.Expr) SelectBuilder {
	for _, expr := range exprs {
		b = b.AddExpr(expr)
	}
	return b
}

func (b SelectBuilder) AddAnys(objects ...any) SelectBuilder {
	var (
		addSqlList = make([]string, len(objects))
		addObjects = make([]any, len(objects))
	)
	for i, object := range objects {
		addSqlList[i] = "?"
		addObjects[i] = object
	}
	return SelectBuilder{
		sqlList: append(b.sqlList, addSqlList...),
		vars:    append(b.vars, addObjects...),
	}
}

func (b SelectBuilder) SqlPattern() string {
	return strings.Join(b.sqlList, ",")
}

func (b SelectBuilder) Vars() []any {
	return b.vars
}
