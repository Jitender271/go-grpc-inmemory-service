package localcache

import (
	"context"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/go-grpc-inmemory-service/internal/log"
	"time"
)

type workerPool struct {
	keyLoadRequestChannel chan interface{}
	config                *Config
	store                 *ristretto.Cache
}

func newWorkerPool(config *Config, store *ristretto.Cache) *workerPool {
	keyLoadReqCh := make(chan interface{}, config.MaxKeyLoadQueueSize)

	wp := &workerPool{
		keyLoadRequestChannel: keyLoadReqCh,
		config:                config,
		store:                 store,
	}

	for i := 0; i < config.Workers; i++ {
		go wp.startWorker()
	}
	return wp
}

func (wp *workerPool) startWorker() {
	for key := range wp.keyLoadRequestChannel {
		fmt.Print("key is : ", key)

		payload, err := runLoader(context.Background(), wp.config.Name, wp.config.LoadTimeout, wp.config.Loader, key)

		if err != nil || payload == nil {
			continue
		}
		value := value{
			payload:   payload,
			createdAt: time.Now(),
		}
		fmt.Print("value is : ", value)

		ok := wp.store.SetWithTTL(key, value, 1, wp.config.TTL)
		if !ok {
			log.Logger.Info("")
		}
	}
}

func (wp *workerPool) enqueue(key interface{}) {
	select {
	case wp.keyLoadRequestChannel <- key:
		return
	default:

	}
}

func (wp *workerPool) destroy() {
	close(wp.keyLoadRequestChannel)
}
