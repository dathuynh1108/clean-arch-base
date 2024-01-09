package usecase

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/repo"
)

type HealthUsecase interface {
	GetHealth() string
}

func NewHealthUsecase(
	healthRepoRouter repo.RepoRouter[*repository.HealthRepository],
) HealthUsecase {
	return &healthUsecase{
		healthRepoRouter: healthRepoRouter,
	}
}

type healthUsecase struct {
	healthRepoRouter repo.RepoRouter[*repository.HealthRepository]
}

func (h *healthUsecase) GetHealth() string {
	return "Hello, World 👋👋👋!"
}
