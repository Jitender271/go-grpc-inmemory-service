package localcache

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"time"
)

type CacheService interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	Delete(ctx context.Context, key interface{})
}

type Cache struct {
	config     *Config
	store      *ristretto.Cache
	workerPool *workerPool
}

func NewCache(config *Config) *Cache {
	store, err := ristretto.NewCache(config.getCacheConfig())
	if err != nil {
		panic(err)
	}
	wp := newWorkerPool(config, store)

	return &Cache{config: config, store: store, workerPool: wp}
}

func (c *Cache) Get(ctx context.Context, key interface{}) (interface{}, error) {

	entry, ok := c.store.Get(key)

	if !ok || entry == nil {
		return c.populateValueFromRemote(ctx, key)
	} else {
		c.refresh(key, entry.(value))
		return entry.(value).payload, nil
	}
}

func (c *Cache) Delete(ctx context.Context, key interface{}) {
	c.store.Del(key)
}

func (c *Cache) refresh(key interface{}, entry value) {
	if time.Now().After(entry.createdAt.Add(c.config.ReloadAfter)) {
		c.workerPool.enqueue(key)
	}
}

func (c *Cache) fillCache(ctx context.Context, key interface{}, dataChan chan interface{}, errChan chan error) {
	payload, err := runLoader(ctx, c.config.Name, c.config.LoadTimeout, c.config.Loader, key)
	var entry value

	if err != nil {
		errChan <- err
	} else {
		entry = value{
			payload:   payload,
			createdAt: time.Now(),
		}
		ok := c.store.SetWithTTL(key, entry, 1, c.config.TTL)
		if !ok {

		}
		dataChan <- entry

	}

}

func (c *Cache) populateValueFromRemote(ctx context.Context, key interface{}) (interface{}, error) {

	dataChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	go c.fillCache(ctx, key, dataChan, errChan)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case dataFromChan := <-dataChan:
		return dataFromChan.(value).payload, nil
	case errFromChan := <-errChan:
		return nil, errFromChan

	}
}
