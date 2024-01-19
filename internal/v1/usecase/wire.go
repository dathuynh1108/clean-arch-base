//go:build wireinject

package usecase

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
	"github.com/google/wire"
)

func ProvideHealthUsecase() HealthUsecase {
	wire.Build(
		repository.ProvideHealthRepo,
		HealthSet,
	)
	return &healthUsecase{}
}
