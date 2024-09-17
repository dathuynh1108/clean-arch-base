//go:build wireinject

package controller

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"

	"github.com/google/wire"
)

func ProvideHealthController() *HealthController {
	wire.Build(
		usecase.ProvideHealthUsecase,
		NewHealthController,
	)
	return &HealthController{}
}

func ProvideErrorController() *ErrorController {
	wire.Build(
		NewErrorController,
	)
	return &ErrorController{}
}
