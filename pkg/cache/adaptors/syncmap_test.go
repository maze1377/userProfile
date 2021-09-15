package adaptors

import (
	"context"
	"testing"
	"userProfile/pkg/cache"

	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {

	syncMapAdaptor := NewSynMapAdaptor()

	if !assert.NotNil(t, syncMapAdaptor) {
		t.Fatal("fail to initialize adaptor")
	}

	key := "test-key"
	value := "test-value"
	ctx := context.Background()
	err := syncMapAdaptor.Set(ctx, key, value)
	assert.NoError(t, err, "fail to set data")
	var resultValue string
	err = syncMapAdaptor.Get(ctx, key, &resultValue)
	assert.NoError(t, err, "fail to get data")
	assert.Equal(t, value, resultValue, "gotten value is not equal to set value")

	err = syncMapAdaptor.Delete(ctx, key)

	if !assert.NoError(t, err) {
		t.Fatal("Error on delete key", err)
	}

	err = syncMapAdaptor.Get(ctx, key, &resultValue)
	assert.Equal(t, err, cache.ErrKeyNotFound)

	err = syncMapAdaptor.Clear(ctx)
	if !assert.NoError(t, err) {
		t.Error("Error on Clear key", err)
	}
}
