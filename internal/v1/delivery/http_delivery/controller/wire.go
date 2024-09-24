//go:build wireinject

package controller

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"

	"github.com/google/wire"
)

func ProvideErrorController() *ErrorController {
	wire.Build(
		NewErrorController,
	)
	return &ErrorController{}
}

func ProvideHealthController() Controller {
	wire.Build(
		usecase.ProvideHealthUsecase,
		HealthSet,
	)
	return &HealthController{}
}
