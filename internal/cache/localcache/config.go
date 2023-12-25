package localcache

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

type Config struct {
	Name                string
	MaxKeys             int
	Loader              Loader
	Workers             int
	MaxKeyLoadQueueSize int
	TTL                 time.Duration
	ReloadAfter         time.Duration
	AsyncLoad           bool
	LoadTimeout         time.Duration
}

func (c *Config) getCacheConfig() *ristretto.Config {
	return &ristretto.Config{
		NumCounters:        int64(c.MaxKeys * 10),
		MaxCost:            int64(c.MaxKeys),
		BufferItems:        64,
		Metrics:            false,
		IgnoreInternalCost: true,
		OnEvict:            func(item *ristretto.Item) {},
		OnReject:           func(item *ristretto.Item) {},
	}
}
