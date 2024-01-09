package dbpool

const (
	AliasMaster  = "master"
	AliasReplica = "replica"
)

type DBAlias string

const (
	// Specific db connections alias here
	DBDefault DBAlias = "default"
)
