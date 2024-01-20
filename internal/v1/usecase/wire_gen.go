// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package usecase

import (
	"github.com/dathuynh1108/clean-arch-base/internal/v1/repository"
)

// Injectors from usecase.go:

func ProvideHealthUsecase() HealthUsecase {
	repoRouter := repository.ProvideHealthRepo()
	usecaseHealthUsecase := NewHealthUsecase(repoRouter)
	return usecaseHealthUsecase
}
