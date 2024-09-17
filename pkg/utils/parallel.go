package utils

import (
	"context"
	"errors"
	"sync"

	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"

	"github.com/panjf2000/ants/v2"
)

func ParallelRun(ctx context.Context, size int, processor ...func(ctx context.Context) error) (err error) {
	// Make size 0 or -1 to use runtime.NumCPU()

	wg := sync.WaitGroup{}
	mu := sync.RWMutex{}
	errs := []error{}

	processCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	onError := func(innerErr error) {
		if innerErr != nil {
			cancel()
			if !errors.Is(innerErr, context.Canceled) {
				mu.Lock()
				defer mu.Unlock()
				errs = append(errs, innerErr)
			}
		}
	}

	pool, err := ants.NewPoolWithFunc(
		size,
		func(i interface{}) {
			defer wg.Done()
			f := i.(func(ctx context.Context) error)
			if err := f(processCtx); err != nil {
				onError(err)
			}
		},
	)
	if err != nil {
		return err
	}
	defer pool.Release()

	for _, p := range processor {
		wg.Add(1)
		if err := pool.Invoke(p); err != nil {
			onError(err)
		}
	}

	wg.Wait()
	if len(errs) > 0 {
		err = comerr.MergeError(errs)
	}
	return
}

func ParallelRunGlobal(ctx context.Context, size int, processor ...func(ctx context.Context) error) (err error) {
	wg := sync.WaitGroup{}
	mu := sync.RWMutex{}
	errs := []error{}

	processCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	onError := func(innerErr error) {
		if innerErr != nil {
			cancel()
			if !errors.Is(innerErr, context.Canceled) {
				mu.Lock()
				defer mu.Unlock()
				errs = append(errs, innerErr)
			}
		}
	}

	for _, p := range processor {
		wg.Add(1)

		wrapProcessor := func(ctx context.Context, processor func(ctx context.Context) error) func() {
			return func() {
				defer wg.Done()

				executeChannel := make(chan error, 1)

				go func() {
					executeChannel <- processor(ctx)
				}()

				select {
				case <-ctx.Done():
					if ctx.Err() != nil {
						onError(ctx.Err())
					}
				case err := <-executeChannel:
					if err != nil {
						onError(err)
					}
				}
			}
		}

		ants.Submit(wrapProcessor(processCtx, p))
	}

	wg.Wait()
	if len(errs) > 0 {
		err = comerr.MergeError(errs)
	}

	return nil
}
