package singleton

import (
	"sync"
)

type Singleton[T any] struct {
	mux      sync.RWMutex
	isLoaded bool
	instance T
	loader   func() T
}

func NewSingleton[T any](loader func() T) *Singleton[T] {
	return &Singleton[T]{
		isLoaded: false,
		loader:   loader,
	}
}

func (s *Singleton[T]) Get() T {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.isLoaded {
		return s.instance
	}

	if !s.isLoaded {
		s.instance = s.loader()
		s.isLoaded = true
	}
	return s.instance
}
