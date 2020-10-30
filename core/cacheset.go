package core

import (
	"github.com/halweg/mini-cache/face"
	"sync"
	"time"
)

type CacheSet struct {
	sync.RWMutex
	items map[string]*CacheItem

	//下次缓存过期检查的间隔时间
	expireInterval time.Duration
	//Timer responsible for triggering cleanup
	//过期清理的定时器
	expireTimer *time.Timer
}

func (set *CacheSet) Add(key string, data interface{}, expire time.Duration) *face.ICacheItem {
	set.Lock()
	defer set.Unlock()

	item := &CacheItem{
		key: key,
		Val: data,
		liveDuration: expire,
		ctime: time.Now(),
	}

	set.items[key] = item
	nextExpire := set.expireInterval

	if expire > 0 && (nextExpire == 0 || expire < nextExpire) {
		set.expireCheck()
	}

	return nil
}

func (set *CacheSet) Get(key string) interface{} {
	set.RLock()
	defer set.RUnlock()
	if val, ok := set.items[key]; ok {
		return val
	}
	return nil
}

func (set *CacheSet) Delete(key string) bool {
	set.Lock()
	defer set.Unlock()

	if _, ok := set.items[key]; ok {
		delete(set.items, key)
	}
	return true
}

func (set *CacheSet) expireCheck() {
	set.Lock()
	//当前有等待执行的定时器任务, 停止该任务并在后续启动新任务
	if set.expireTimer == nil {
		set.expireTimer.Stop()
	}

	now := time.Now()
	nextCheckTime := 0
	for k, v := range set.items {
		// item which will never be overdue
		if v.liveDuration == 0 {
			continue
		}

		//判断键是否已经过期了，过期则删除，未过期的话则找出下一次最近的过期时间
		if now.Sub(v.ctime) >= v.liveDuration {
			delete(set.items, k)
		} else {
			if nextCheckTime == 0 || nextCheckTime > int(v.liveDuration - now.Sub(v.ctime)) {
				nextCheckTime = int(v.liveDuration - now.Sub(v.ctime))
			}
		}
	}

	//如果缓存中还有kv存在，设置定时器定期evict过期键
	if nextCheckTime > 0 {
		set.expireTimer = time.AfterFunc(time.Duration(nextCheckTime), func() {
			set.expireCheck()
		})
	}
	set.Unlock()
}
