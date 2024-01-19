package usecase

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/repo"
	"github.com/google/wire"
)

var (
	HealthSet = wire.NewSet(
		NewHealthUsecase,
		wire.Bind(new(HealthUsecase), new(*healthUsecase)),
	)
)

type HealthUsecase interface {
	GetHealth() string
}

func NewHealthUsecase(
	healthRepoRouter repo.RepoRouter[repository.HealthRepository],
) *healthUsecase {
	return &healthUsecase{
		healthRepoRouter: healthRepoRouter,
	}
}

type healthUsecase struct {
	healthRepoRouter repo.RepoRouter[repository.HealthRepository]
}

func (h *healthUsecase) GetHealth() string {
	return "Hello, World 👋👋👋!"
}
