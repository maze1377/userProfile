package multilayercache

import (
	"context"
	"reflect"
	"sync"
	"userProfile/pkg/cache"
	"userProfile/pkg/cache/adaptors"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type multilayerCache struct {
	layers []cache.Layer
}

type layerOperator func(layer cache.Layer) (interface{}, error)

func New(layers ...cache.Layer) cache.Layer {
	return &multilayerCache{
		layers: layers,
	}
}

func (mc *multilayerCache) Get(ctx context.Context, key string, reference interface{}) error {
	value, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		// todo create new item from reference
		items := reflect.ValueOf(reference).Elem().Interface()
		err := layer.Get(ctx, key, &items)
		return items, err
	})
	if err == nil {
		adaptors.DeepCopy(value, reference)
	}
	return err
}

func (mc *multilayerCache) Delete(ctx context.Context, key string) error {
	_, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		err := layer.Delete(ctx, key)
		return nil, err
	})

	return err
}

func (mc *multilayerCache) Set(ctx context.Context, key string, value interface{}) error {
	_, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		err := layer.Set(ctx, key, value)
		return nil, err
	})

	return err
}

func (mc *multilayerCache) Clear(ctx context.Context) error {
	_, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		err := layer.Clear(ctx)
		return nil, err
	})

	return err
}

func (mc *multilayerCache) wrapAllErrors(errChannel <-chan error) error {
	var allErrors error
	for err := range errChannel {
		if allErrors == nil {
			allErrors = err
		} else {
			allErrors = errors.Wrap(err, allErrors.Error())
		}
	}

	return allErrors
}

func (mc *multilayerCache) performOperation(operator layerOperator) (interface{}, error) {
	errChannel := make(chan error, len(mc.layers))
	resultChannel := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(mc.layers))

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errChannel)
		close(resultChannel)
	}()

	for _, cacheLayer := range mc.layers {
		go func(layer cache.Layer) {
			defer wg.Done()

			value, err := operator(layer)
			if err != nil {
				errChannel <- err
				return
			}

			select {
			case resultChannel <- value:
			default:
			}
		}(cacheLayer)
	}

	select {
	case value := <-resultChannel:
		go func() {
			for err := range errChannel {
				logrus.WithError(err).Warning("failed to get value from multilayer cache")
			}
		}()

		return value, nil

	case <-done:
		return nil, mc.wrapAllErrors(errChannel)
	}
}
