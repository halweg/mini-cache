package core

import (
	"sync"
	"time"
)

type CachePool struct {
	sync.RWMutex
	items map[string]interface{}
}

func (ca *CachePool) Put(key string, data map[string]interface{}) {
	ca.Lock()
	defer ca.Unlock()
	ca.items[key] = data
}

func (ca *CachePool) Get(key string) interface{}  {
	ca.RLock()

	if v, ok := ca.items[key]; ok {
		if val, ok := v.(map[string]interface{}); ok {
			ca.RUnlock()
			if time.Now().Unix() < val["expires"].(int64) {
				return val["content"]
			} else {
				ca.Lock()
				delete(ca.items, key)
				ca.Unlock()
				return nil
			}
		}
	}

	ca.RUnlock()
	return nil
}

func (ca *CachePool) Del(key string) bool {
	ca.Lock()
	defer ca.Unlock()

	if _, ok := ca.items[key]; ok {
		delete(ca.items, key)
		return true
	}

	return false
}