package singleton

import (
	"fmt"
	"sync"
)

type Singleton[T any] struct {
	instance T
	loader   func() T
	once     sync.Once
}

func NewSingleton[T any](loader func() T, preload bool) *Singleton[T] {
	s := &Singleton[T]{
		loader: loader,
	}
	if preload {
		_ = s.Get()
	}
	return s
}

func (s *Singleton[T]) Get() T {
	s.once.Do(func() {
		s.instance = s.loader()
	})

	return s.instance
}

type SingletonMap[T any] struct {
	mux         sync.Mutex
	instanceMap map[string]T
	loader      func(key fmt.Stringer) T
}

func NewSingletonMap[T any](loader func(key fmt.Stringer) T) *SingletonMap[T] {
	return &SingletonMap[T]{
		loader:      loader,
		instanceMap: make(map[string]T),
	}
}

func (s *SingletonMap[T]) Get(key fmt.Stringer) T {
	textKey := key.String()
	if value, ok := s.instanceMap[textKey]; ok {
		return value
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	value, ok := s.instanceMap[textKey]
	if !ok {
		value = s.loader(key)
		s.instanceMap[textKey] = value
	}
	return value
}
