package cupcakecache

import (
	"errors"
	"sync"
)

// GetterFunc is for user to define how to get data from database
// when the application fails to hit the cache with given key
type GetterFunc func(key string) ([]byte, error)

// GetterFunc implements the interface Getter
func (getter GetterFunc) Get(key string) ([]byte, error) {
	return getter(key)
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type Group struct {
	name   string
	cache  *cache
	getter Getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, size int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()

	group := &Group{
		name:   name,
		cache:  &cache{size: size},
		getter: getter,
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("got empty key")
	}
	// First try to get data from cache
	bv, ok := g.cache.get(key)
	if ok {
		return bv, nil
	}
	// TODO: If fail to hit the cache on a specific node, then search from another node

	// If fail to hit the cache, then get from database and then save to cache
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	val := make([]byte, len(bytes))
	// Add data to cache
	g.cache.add(key, ByteView{val})
	copy(val, bytes)
	return ByteView{val}, nil
}
