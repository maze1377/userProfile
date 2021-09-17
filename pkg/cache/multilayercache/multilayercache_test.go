package multilayercache

import (
	"context"
	"testing"
	"userProfile/pkg/cache"
	"userProfile/pkg/cache/adaptors"

	"golang.org/x/xerrors"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	multiLayerCache := New(adaptors.NewSynMapAdaptor(), adaptors.NewSynMapAdaptor())

	key := "test-key"
	value := "test-value"
	ctx := context.Background()

	err := multiLayerCache.Set(ctx, key, value)
	assert.NoError(t, err, "fail to set data")
	var resultValue string
	err = multiLayerCache.Get(ctx, key, &resultValue)
	assert.NoError(t, err, "fail to get data")

	assert.Equal(t, value, resultValue, "gotten value is not equal to set value")

	err = multiLayerCache.Delete(ctx, key)

	if !assert.NoError(t, err) {
		t.Fatal("Error on delete key", err)
	}

	err = multiLayerCache.Get(ctx, key, &resultValue)

	if !xerrors.Is(err, cache.ErrKeyNotFound) {
		t.Error("find deleted key!", err)
	}

	err = multiLayerCache.Clear(ctx)
	if !assert.NoError(t, err) {
		t.Error("Error on Clear key", err)
	}
}
