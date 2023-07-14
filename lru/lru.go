package lru

import (
	"container/list"
)

type LRUCache struct {
	maxBytes   int64
	curBytes   int64
	linkedList *list.List
	data       map[string]*list.Element
	// Optional function that will be called when a key-value pair is removed
	onEvited func(key string, val Value)
}

type entry struct {
	key string
	// interface embed in struct
	// value can be any base type or struct as long as it has Len() method
	value Value
}

type Value interface {
	Len() int
}

// Cache does not have size limit when maxBytes is set to zero
func New(maxBytes int64, onEvited func(key string, val Value)) *LRUCache {
	return &LRUCache{
		maxBytes:   maxBytes,
		curBytes:   0,
		linkedList: list.New(),
		data:       make(map[string]*list.Element),
		onEvited:   onEvited,
	}
}

func (cache *LRUCache) Get(key string) (Value, bool) {
	if elem, ok := cache.data[key]; ok {
		// When element be accessed, move it to the end of the list
		cache.linkedList.MoveToBack(elem)
		pair := elem.Value.(*entry)
		return pair.value, ok
	}
	return nil, false
}

// Remove the oldest key-value pair from cache, remove the head of the list
func (cache *LRUCache) RemoveOldest() {
	if elem := cache.linkedList.Front(); elem != nil {
		key := elem.Value.(*entry).key
		// Remove element from head
		cache.linkedList.Remove(elem)
		// Remove data from map
		delete(cache.data, key)
		// Update current size of the cache
		pair := elem.Value.(*entry)
		cache.curBytes -= (int64(len(key)) + int64(pair.value.Len()))
		// If OnEvited exist, then call with the value of the deleted pair
		if cache.onEvited != nil {
			cache.onEvited(key, pair.value)
		}
	}
}

func (cache *LRUCache) Add(key string, val Value) {
	if elem, ok := cache.data[key]; ok {
		// If key exists, update the value and move it to the end of the list
		entry := elem.Value.(*entry)
		// Update current size of the cache
		cache.curBytes = cache.curBytes - int64(entry.value.Len()) + int64(val.Len())
		entry.value = val
		cache.linkedList.MoveToBack(elem)
	} else {
		//If key doesn't exist, add new key-value pair to map and end of the list
		elem = cache.linkedList.PushBack(&entry{key: key, value: val})
		cache.data[key] = elem
		// Update current size of the cache
		cache.curBytes += (int64(len(key)) + int64(val.Len()))
	}

	// If current size of the cache is larger then the maximum size of cache
	// then remove the oldest record(s)
	for cache.maxBytes != 0 && cache.curBytes > cache.maxBytes {
		cache.RemoveOldest()
	}
}

func (c *LRUCache) Len() int {
	return c.linkedList.Len()
}
