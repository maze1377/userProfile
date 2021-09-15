package adaptors

import (
	"context"
	"sync"
	"userProfile/pkg/cache"
)

type syncMap struct {
	sync.RWMutex
	data sync.Map
}

func NewSynMapAdaptor() cache.Layer {
	return &syncMap{
		data: sync.Map{},
	}
}
func (syncInstance *syncMap) Set(_ context.Context, key string, value interface{}) error {
	syncInstance.RLock()
	defer syncInstance.RUnlock()
	syncInstance.data.Store(key, value)
	return nil
}

func (syncInstance *syncMap) Get(_ context.Context, key string, reference interface{}) error {
	syncInstance.RLock()
	defer syncInstance.RUnlock()
	value, ok := syncInstance.data.Load(key)
	if !ok {
		return cache.ErrKeyNotFound
	}
	DeepCopy(value, reference)
	return nil
}

func (syncInstance *syncMap) Delete(_ context.Context, key string) error {
	syncInstance.RLock()
	defer syncInstance.RUnlock()
	syncInstance.data.Delete(key)
	return nil
}

func (syncInstance *syncMap) Clear(_ context.Context) error {
	syncInstance.Lock()
	defer syncInstance.Unlock()
	syncInstance.data = sync.Map{}
	return nil
}
