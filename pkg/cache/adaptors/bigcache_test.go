package adaptors

import (
	"context"
	"testing"
	"time"
	"userProfile/pkg/cache"

	"github.com/allegro/bigcache/v2"
	"github.com/stretchr/testify/assert"
)

func TestBigCacheSetGet(t *testing.T) {
	bigCacheInstance, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             1,
		LifeWindow:         2 * time.Minute,
		MaxEntriesInWindow: 128,
		CleanWindow:        1 * time.Minute,
		HardMaxCacheSize:   128,
		Verbose:            true,
	})

	if !assert.NoError(t, err) {
		t.Fatal("fail to initialize big cache")
	}

	bigCacheAdaptor := NewBigCacheAdaptor(bigCacheInstance)

	if !assert.NotNil(t, bigCacheAdaptor) {
		t.Fatal("fail to initialize adaptor")
	}

	key := "test-key"
	value := "test-value"
	ctx := context.Background()
	err = bigCacheAdaptor.Set(ctx, key, value)
	assert.NoError(t, err, "fail to set data")
	var resultValue string
	err = bigCacheAdaptor.Get(ctx, key, &resultValue)
	assert.NoError(t, err, "fail to get data")
	assert.Equal(t, value, resultValue, "gotten value is not equal to set value")

	err = bigCacheAdaptor.Delete(ctx, key)

	if !assert.NoError(t, err) {
		t.Fatal("Error on delete key", err)
	}
	err = bigCacheAdaptor.Get(ctx, key, &resultValue)
	assert.Equal(t, err, cache.ErrKeyNotFound)

	err = bigCacheAdaptor.Clear(ctx)
	if !assert.NoError(t, err) {
		t.Error("Error on Clear key", err)
	}
}
