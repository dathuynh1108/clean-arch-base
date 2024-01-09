package dbpool

import "strings"

func BuildAlias(alias, instanceAlias string) string {
	return alias + "." + instanceAlias
}

func ParseAlias(alias string) (string, string) {
	parts := strings.Split(alias, ".")
	return parts[0], parts[1]
}
