package sqlquery

import (
	"gorm.io/gorm/clause"
)

// ExOrderBy wraps clause.OrderBy and support DESC for Expression.
// When MergeClause, it respects the clause's expression and chain with its Expression
// by CommaChainExpression.
type ExOrderBy struct {
	clause.OrderBy
	IsExpressionDesc bool
}

func (ob ExOrderBy) Build(builder clause.Builder) {
	if ob.Expression == nil {
		ob.OrderBy.Build(builder)
	} else {
		ob.Expression.Build(builder)
		if ob.IsExpressionDesc {
			_, _ = builder.WriteString(" DESC")
		}
	}
}

func (ob ExOrderBy) MergeClause(c *clause.Clause) {
	expression := c.Expression
	ob.OrderBy.MergeClause(c)
	if expression == nil {
		c.Expression = ob
	} else {
		c.Expression = CommaChainExpression{
			Expressions: []clause.Expression{expression, ob},
		}
	}
}
