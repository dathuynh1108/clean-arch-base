package sqlquery

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type fulltextSearchClause struct {
	Columns []clause.Column
	Text    string
}

func (m *fulltextSearchClause) Build(builder clause.Builder) {
	stmt, ok := builder.(*gorm.Statement)
	if !ok {
		m.BuildPostgres(builder)
		return
	}
	switch stmt.DB.Name() {
	case "postgres":
		m.BuildPostgres(builder)
	case "mysql":
		m.BuildMySQL(builder)
	default:
		m.BuildPostgres(builder)
	}
}

func (m *fulltextSearchClause) BuildMySQL(builder clause.Builder) {
	builder.WriteString("MATCH (")
	if len(m.Columns) > 0 {
		for _, col := range m.Columns[:len(m.Columns)-1] {
			builder.WriteQuoted(col)
			builder.WriteByte(',')
		}
		builder.WriteQuoted(m.Columns[len(m.Columns)-1])
	}
	builder.WriteString(") AGAINST (")
	builder.AddVar(builder, m.Text)
	builder.WriteString(" IN BOOLEAN MODE)")
}

func (m *fulltextSearchClause) BuildPostgres(builder clause.Builder) {
	builder.WriteString("to_tsvector('english', ")
	if len(m.Columns) > 0 {
		builder.WriteQuoted(m.Columns[0])
		for _, col := range m.Columns[1:] {
			builder.WriteString(" || ' ' || ")
			builder.WriteQuoted(col)
		}
	}
	builder.WriteString(") @@ plainto_tsquery(")
	builder.AddVar(builder, m.Text)
	builder.WriteString(")")
}
