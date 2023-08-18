package usecase

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
	repositoryrouter "github.com/dathuynh1108/clean-arch-base/internal/v1/repository_router"
)

type HealthUsecase interface {
	GetHealth() string
}

func NewHealthUsecase(
	healthRepoRouter repositoryrouter.RepositoryRouter[*repository.HealthRepository],
) HealthUsecase {
	return &healthUsecase{
		healthRepoRouter: healthRepoRouter,
	}
}

type healthUsecase struct {
	healthRepoRouter repositoryrouter.RepositoryRouter[*repository.HealthRepository]
}

func (h *healthUsecase) GetHealth() string {
	return "Hello, World 👋👋👋!"
}
