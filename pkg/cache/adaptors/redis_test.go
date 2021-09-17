package adaptors

import (
	"context"
	"reflect"
	"testing"
	"userProfile/pkg/cache"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func newTestRedis() string {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return mr.Addr()
}

func Test_Redis(t *testing.T) {
	masterServer := newTestRedis()
	ri, err := NewRedisAdaptor(&InstanceOptions{
		Address: InstanceOptionsAddress{
			Master: masterServer,
			Replicas: []string{
				masterServer,
				masterServer,
			},
		},
	})

	if err != nil {
		t.Error("Error redis connection", err)
	}

	ctx := context.Background()

	key := "test-key"
	value := "test-value"

	// Test Set value
	err = ri.Set(ctx, key, value)
	if err != nil {
		t.Error("Error on Set in Redis", err)
	}

	// Test Get value
	var resultValue string
	err = ri.Get(ctx, key, &resultValue)
	if !assert.NoError(t, err) {
		t.Error("Error on Get in Redis", err)
	}

	// Compare actual data vs cached data
	if eq := reflect.DeepEqual(resultValue, value); !eq {
		t.Error("Expected", value, "Got", resultValue)
	}

	// Getting non-existing value
	err = ri.Get(ctx, "non_existing_key", &resultValue)
	if !assert.Equal(t, err, cache.ErrKeyNotFound) {
		t.Error("Error in Getting non-existing value, Expected some error")
	}

	// delete key from cache
	err = ri.Delete(ctx, key)
	if !assert.NoError(t, err) {
		t.Error("Error on delete key", err)
	}

	err = ri.Get(ctx, key, &resultValue)
	assert.Equal(t, err, cache.ErrKeyNotFound)

	err = ri.Clear(ctx)
	if !assert.NoError(t, err) {
		t.Error("Error on delete key", err)
	}
}
