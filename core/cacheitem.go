package core

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex
	key string
	Val interface{}

	//存活时间
	liveDuration time.Duration

	//创建时间
	ctime time.Time
}
