package repositoryrouter

import (
	"sync"
	"sync/atomic"
)

type RepositoryRouter[Repo any] interface {
	GetRepo(alias string) Repo
	SetRepo(alias string, repo Repo)
}

func NewRepositoryRouter[Repo any]() RepositoryRouter[Repo] {
	return &repositoryRouter[Repo]{
		mapper: map[string][]Repo{},
		score:  map[string]*atomic.Uint32{},
	}
}

type repositoryRouter[Repo any] struct {
	mux    sync.RWMutex
	mapper map[string][]Repo
	score  map[string]*atomic.Uint32
}

func (r *repositoryRouter[Repo]) GetRepo(alias string) Repo {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.mapper[alias][r.score[alias].Add(1)%uint32(len(r.mapper[alias]))]
}

func (r *repositoryRouter[Repo]) GetRepoForce(alias string, index uint32) Repo {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.mapper[alias][index]
}

func (r *repositoryRouter[Repo]) SetRepo(alias string, repo Repo) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.mapper[alias] = append(r.mapper[alias], repo)
	r.score[alias] = &atomic.Uint32{}
}
