package face

import "time"

type ICacheSet interface {

	Get(key string) interface{}

	Delete (key string) bool

	Add(key string, data interface{}, expire time.Duration) *ICacheItem

}
