package core

import (
    "container/list"
    "github.com/halweg/mini-cache/face"
)

type Cache struct {
    maxBytes int64
    nBytes int64

    ll *list.List
    cache map[string]*list.Element

    OnEvicted func(key string, value face.IValue)
}

type entry struct {
    key string
    value face.IValue
}

func NewValue(maxBytes int64, OnEvicted func(key string, value face.IValue)) *Cache {
    return &Cache{
        maxBytes:  maxBytes,
        ll:        list.New(),
        cache:     make(map[string]*list.Element),
        OnEvicted: OnEvicted,
    }
}




