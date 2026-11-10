package ctxlog

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentAddAndGet(t *testing.T) {
	ctx := New(context.Background())

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			Add(ctx, fmt.Sprintf("key_%d", i), i)
		}(i)
		go func() {
			defer wg.Done()
			for range Get(ctx) {
			}
		}()
	}
	wg.Wait()

	got := Get(ctx)
	if len(got) != 50 {
		t.Fatalf("expected 50 fields, got %d", len(got))
	}
}

func TestGetReturnsCopy(t *testing.T) {
	ctx := New(context.Background())
	Add(ctx, "a", 1)

	got := Get(ctx)
	got["b"] = 2

	if _, ok := Get(ctx)["b"]; ok {
		t.Fatal("mutating Get result must not affect stored fields")
	}
}

func TestNewReusesFields(t *testing.T) {
	ctx := New(context.Background())
	child := New(context.WithValue(ctx, struct{}{}, "x"))
	Add(child, "from child", true)

	if v, ok := Get(ctx)["from_child"]; !ok || v != true {
		t.Fatalf("expected child Add visible in parent, got %v", Get(ctx))
	}
}

func TestNilAndMissingContext(t *testing.T) {
	//lint:ignore SA1012 nil context intentionally tested
	Add(nil, "k", "v")
	//lint:ignore SA1012 nil context intentionally tested
	if Get(nil) != nil {
		t.Fatal("Get(nil) should be nil")
	}
	Add(context.Background(), "k", "v")
	if Get(context.Background()) != nil {
		t.Fatal("Get without New should be nil")
	}
}
