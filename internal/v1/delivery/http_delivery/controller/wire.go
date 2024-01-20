//go:build wireinject

package controller

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/usecase"
	"github.com/google/wire"
)

func ProvideHealthController() *HealthControler {
	wire.Build(
		usecase.ProvideHealthUsecase,
		NewHealthController,
	)
	return &HealthControler{}
}

func ProvideWSController() *WSController {
	wire.Build(
		NewWSController,
	)
	return &WSController{}
}

func ProvideErrorController() *ErrorControler {
	wire.Build(
		NewErrorController,
	)
	return &ErrorControler{}
}
