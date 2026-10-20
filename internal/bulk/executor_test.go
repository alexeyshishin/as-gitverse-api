package bulk

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRun_PartialFailurePreservesOrder(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	results := Run(context.Background(), items, 2, func(ctx context.Context, i int) error {
		if i%2 == 0 {
			return errors.New("even fail")
		}
		return nil
	})

	if len(results) != len(items) {
		t.Fatalf("got %d results, want %d", len(results), len(items))
	}
	for idx, r := range results {
		if r.Item != items[idx] {
			t.Fatalf("result %d: item %d, want %d", idx, r.Item, items[idx])
		}
		wantFail := items[idx]%2 == 0
		if (r.Err != nil) != wantFail {
			t.Fatalf("result %d: err=%v, wantFail=%v", idx, r.Err, wantFail)
		}
	}
}

func TestRun_ContextCancelStopsRemaining(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{}, 1)

	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	results := Run(ctx, items, 1, func(ctx context.Context, i int) error {
		select {
		case started <- struct{}{}:
		default:
		}
		if i == 1 {
			cancel()
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	})

	canceled := 0
	for _, r := range results {
		if errors.Is(r.Err, context.Canceled) {
			canceled++
		}
	}
	if canceled == 0 {
		t.Fatalf("expected at least one canceled result")
	}
}

func TestRun_ConcurrencyFloor(t *testing.T) {
	results := Run(context.Background(), []int{1}, 0, func(ctx context.Context, i int) error { return nil })
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("got %+v", results)
	}
}
