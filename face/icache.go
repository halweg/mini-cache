package face

type ICache interface {
	Put(key string, data map[string]interface{})

	Get(key string) interface{}

	Del(key string) bool
}
