package utils

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"log"
)

type ActionCache struct {
	data    *lru.Cache[interface{}, interface{}]
	getFunc func(key interface{}) interface{}
}

func NewActionCache(cacheSize int) *ActionCache {
	cache, err := lru.New[interface{}, interface{}](cacheSize)
	if err != nil {
		log.Fatal(err)
	}
	return &ActionCache{
		data: cache,
	}
}

func (ac *ActionCache) Get(key interface{}) interface{} {
	value, ok := ac.data.Get(key)
	if !ok {
		value = ac.getFunc(key)
		if value != nil {
			ac.data.Add(key, value)
			return value
		}
		return nil
	}
	return value
}

func (ac *ActionCache) RegisterGetFunc(getFunc func(key interface{}) interface{}) {
	ac.getFunc = getFunc
}

func (ac *ActionCache) Delete(key interface{}) {
	_, ok := ac.data.Get(key)
	if !ok {
		return
	}
	ac.data.Remove(key)
}
