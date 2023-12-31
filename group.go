package cupcake_cache

import (
	"errors"
	"sync"

	"github.com/lz-nsc/cupcake_cache/log"
)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
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
	name    string
	cache   *cache
	getter  Getter
	remote  Server
	callMgr *CallManager
}

func NewGroup(name string, size int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	mu.Lock()
	defer mu.Unlock()

	group := &Group{
		name:    name,
		cache:   &cache{size: size},
		getter:  getter,
		callMgr: NewCallManger(),
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) SetRemote(remote Server) {
	if g.remote != nil {
		panic("remote handler already exists")
	}
	g.remote = remote
}

func (g *Group) Get(key string) (ByteView, error) {
	if g.callMgr == nil {
		g.callMgr = NewCallManger()
	}
	val, err := g.callMgr.Do(key, func() (interface{}, error) {
		if key == "" {
			return ByteView{}, errors.New("got empty key")
		}
		// First try to get data from cache
		bv, ok := g.cache.get(key)
		if ok {
			log.Debugf("successfully hit local cache, key: %s", key)
			return bv, nil
		}

		return g.getFromRemote(key)
	})

	return val.(ByteView), err

}
func (g *Group) getFromRemote(key string) (ByteView, error) {
	// Get from remote cache server
	bytes, err := g.getFromRemoteCache(key)
	if err != nil {
		log.Errorf("failed to get from remote cache, err: %s", err.Error())
	}

	if bytes == nil {
		log.Debugf("get record from database, key: %s", key)
		// If fail to hit the cache, then get from database and then save to cache
		bytes, err = g.getter.Get(key)
		if err != nil {
			return ByteView{}, err
		}
		log.Debugf("successfully get record from database, key: %s", key)

		val := make([]byte, len(bytes))
		// Add data to local cache
		g.cache.add(key, ByteView{val})
		copy(val, bytes)
	}

	return ByteView{bytes: bytes}, nil
}

func (g *Group) getFromRemoteCache(key string) ([]byte, error) {
	if g.remote == nil {
		return nil, nil
	}
	log.Debugf("get record from remote cache, key: %s", key)

	bytes, err := g.remote.RemoteGet(g.name, key)
	if err != nil {
		return nil, err
	}

	log.Debugf("successfully get record from remote cache, key: %s", key)
	return bytes, nil
}
