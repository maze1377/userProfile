package multilayercache

import (
	"context"
	"testing"
	"userProfile/pkg/cache/adaptors"

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
	// todo assert error wrapped by key not NotFound

	err = multiLayerCache.Clear(ctx)
	if !assert.NoError(t, err) {
		t.Error("Error on Clear key", err)
	}
}
