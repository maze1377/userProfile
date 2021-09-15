package cache

import "context"

type Layer interface {
	Get(ctx context.Context, key string, reference interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
