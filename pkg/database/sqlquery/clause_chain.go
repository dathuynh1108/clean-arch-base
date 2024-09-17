package sqlquery

import (
	"gorm.io/gorm/clause"
)

type CommaChainExpression struct {
	Expressions []clause.Expression
}

func (ce CommaChainExpression) Build(builder clause.Builder) {
	if len(ce.Expressions) == 0 {
		_, _ = builder.WriteString("0") // Failsafe
		return
	}
	lastIdx := len(ce.Expressions) - 1
	for i, e := range ce.Expressions {
		e.Build(builder)
		if i < lastIdx {
			_ = builder.WriteByte(',')
		}
	}
}

func (ce *CommaChainExpression) AddExpression(expr clause.Expression) {
	ce.Expressions = append(ce.Expressions, expr)
}
