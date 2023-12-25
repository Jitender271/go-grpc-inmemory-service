package localcache

import "time"

type value struct {
	payload   interface{}
	createdAt time.Time
}
