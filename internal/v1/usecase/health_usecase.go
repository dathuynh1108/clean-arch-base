package usecase

import (
	"context"

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
	GetHealth(ctx context.Context) string
}

func NewHealthUsecase(
	healthRepoRouter repo.RepoRouter[repository.HealthRepo],
) *healthUsecase {
	return &healthUsecase{
		healthRepoRouter: healthRepoRouter,
	}
}

type healthUsecase struct {
	healthRepoRouter repo.RepoRouter[repository.HealthRepo]
}

func (h *healthUsecase) GetHealth(ctx context.Context) string {
	return "Hello, World 👋👋👋!"
}
