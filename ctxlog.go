// Package ctxlog provides a way to add context fields to a context.
package ctxlog

import (
	"context"
	"strings"
	"sync"
)

type ctxLogKey struct{}

type fields struct {
	mu sync.RWMutex
	m  map[string]any
}

var ctxLogKeyFields = ctxLogKey{}

func New(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := ctx.Value(ctxLogKeyFields).(*fields); ok {
		return ctx
	}

	return context.WithValue(ctx, ctxLogKeyFields, &fields{m: map[string]any{}})
}

// Add sets a field on the context. It is safe for concurrent use.
func Add(ctx context.Context, key string, value any) {
	f := fromContext(ctx)
	if f == nil {
		return
	}

	key = strings.ReplaceAll(key, " ", "_")

	f.mu.Lock()
	f.m[key] = value
	f.mu.Unlock()
}

// Get returns a snapshot copy of the context fields. It is safe for concurrent use.
func Get(ctx context.Context) map[string]any {
	f := fromContext(ctx)
	if f == nil {
		return nil
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	out := make(map[string]any, len(f.m))
	for k, v := range f.m {
		out[k] = v
	}
	return out
}

func fromContext(ctx context.Context) *fields {
	if ctx == nil {
		return nil
	}
	f, _ := ctx.Value(ctxLogKeyFields).(*fields)
	return f
}
