package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	n := len(in)
	if n == 0 {
		return []Result[R]{}, nil
	}
	if workers > n {
		workers = n
	}

	results := make([]Result[R], n)
	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case i, ok := <-jobs:
					if !ok {
						return
					}
					value, err := fn(ctx, in[i])
					if err != nil {
						results[i] = Result[R]{Err: err}
						continue
					}
					results[i] = Result[R]{Value: value}
				}
			}
		}()
	}

	next := 0
send:
	for next < n {
		select {
		case <-ctx.Done():
			break send
		case jobs <- next:
			next++
		}
	}
	close(jobs)
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
