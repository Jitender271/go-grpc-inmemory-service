package localcache

import (
	"context"
	"time"
)

type Loader func(ctx context.Context, key interface{}) (interface{}, error)

func runLoader(ctx context.Context, cacheName string, timeout time.Duration, loader Loader, key interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	payloadChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	go func() {
		payload, err := loader(ctx, key)

		if err != nil {
			errChan <- err
		} else {
			payloadChan <- payload
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case payload := <-payloadChan:
		return payload, nil
	}

}
