package usecase

import "github.com/dathuynh1108/clean-arch-base/pkg/singleton"

var (
	HealthUsecaseSingleton = singleton.NewSingleton(ProvideHealthUsecase, false)
)
