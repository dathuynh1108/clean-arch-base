package dbpool

import (
	"sync"
	"sync/atomic"
)

type DBPool[DB any] interface {
	GetMaster(alias DBAlias) DB
	GetReplica(alias DBAlias) DB
	GetDB(alias string) DB
	GetDBForce(alias string, index uint32) DB
	SetDB(alias string, repo DB)
}

func NewDBPool[DB any]() DBPool[DB] {
	return &dbPool[DB]{
		mapper: map[string][]DB{},
		score:  map[string]*atomic.Uint32{},
	}
}

type dbPool[DB any] struct {
	mux    sync.RWMutex
	mapper map[string][]DB
	score  map[string]*atomic.Uint32
}

func (r *dbPool[DB]) GetMaster(alias DBAlias) DB {
	return r.GetDB(BuildAlias(alias, AliasMaster))
}

func (r *dbPool[DB]) GetReplica(alias DBAlias) DB {
	return r.GetDB(BuildAlias(alias, AliasMaster))
}

func (r *dbPool[DB]) GetDB(alias string) DB {
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.mapper[alias][r.score[alias].Add(1)%uint32(len(r.mapper[alias]))]
}

func (r *dbPool[DB]) GetDBForce(alias string, index uint32) DB {
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.mapper[alias][index]
}

func (r *dbPool[DB]) SetDB(alias string, db DB) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.mapper[alias] = append(r.mapper[alias], db)
	r.score[alias] = &atomic.Uint32{}
}
