package bulk

import (
	"context"
	"sync"
)

type Result[I any] struct {
	Item I
	Err  error
}

func Run[I any](ctx context.Context, items []I, concurrency int, fn func(context.Context, I) error) []Result[I] {
	if concurrency < 1 {
		concurrency = 1
	}

	results := make([]Result[I], len(items))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(i int, item I) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = Result[I]{Item: item, Err: ctx.Err()}
				return
			}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				results[i] = Result[I]{Item: item, Err: ctx.Err()}
				return
			}

			results[i] = Result[I]{Item: item, Err: fn(ctx, item)}
		}(i, item)
	}

	wg.Wait()
	return results
}
