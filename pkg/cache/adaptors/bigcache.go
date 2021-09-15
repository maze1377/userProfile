package adaptors

import (
	"context"
	"userProfile/pkg/cache"

	"github.com/allegro/bigcache/v2"
)

type bigCache struct {
	ins *bigcache.BigCache
}

func NewBigCacheAdaptor(instance *bigcache.BigCache) cache.Layer {
	ins := &bigCache{
		ins: instance,
	}
	return ins
}

func (c *bigCache) Set(_ context.Context, key string, value interface{}) error {
	encodedValue, err := encode(value)
	if err != nil {
		return cache.ErrEncode
	}
	err = c.ins.Set(key, encodedValue)
	return err
}

func (c *bigCache) Get(_ context.Context, key string, reference interface{}) error {
	rawData, err := c.ins.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return cache.ErrKeyNotFound
	} else if err != nil {
		return err
	}
	err = decode(rawData, reference)

	if err != nil {
		return cache.ErrDecode
	}
	return nil
}

func (c *bigCache) Delete(_ context.Context, key string) error {
	err := c.ins.Delete(key)
	return err
}

func (c *bigCache) Clear(_ context.Context) error {
	err := c.ins.Reset()
	return err
}
