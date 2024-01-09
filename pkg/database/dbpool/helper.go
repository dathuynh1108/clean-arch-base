package dbpool

import (
	"fmt"
	"strings"
)

func BuildAlias(alias DBAlias, instanceAlias string) string {
	return fmt.Sprintf("%s.%s", alias, instanceAlias)
}

func ParseAlias(alias string) (DBAlias, string) {
	parts := strings.Split(alias, ".")
	return DBAlias(parts[0]), parts[1]
}
