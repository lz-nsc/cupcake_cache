package cupcake_cache

import (
	"fmt"
	"strings"
	"sync"

	hash "github.com/lz-nsc/cupcake_cache/Hash"
)

const (
	defaultGrpcReplicas = 4
)

type Server interface {
	Run() error
	RemoteGet(group string, key string) ([]byte, error)
}
type generator func(string, string, *CacheServer) Server

var serverMap = map[string]generator{}

type CacheServer struct {
	server    Server
	name      string
	addr      string
	remoteMap map[string]string
	mu        sync.Mutex
	peerMgr   *hash.Manager
}

func NewServer(proto string, name string, addr string, hashFn hash.HashFunc) (*CacheServer, error) {
	cs := &CacheServer{
		name:      name,
		addr:      addr,
		mu:        sync.Mutex{},
		remoteMap: map[string]string{},
		peerMgr:   hash.NewManager(defaultGrpcReplicas, hashFn),
	}
	fn := serverMap[proto]
	if fn == nil {
		return nil, fmt.Errorf("invalid protocol, got: %s", proto)
	}
	cs.server = fn(name, addr, cs)
	return cs, nil
}

func (cs *CacheServer) Run() error {
	return cs.server.Run()
}

func (cs *CacheServer) AddGroup(name string, size int64, getter Getter) *Group {
	group := NewGroup(name, size, getter)
	group.SetRemote(cs.server)
	return group
}

func (cs *CacheServer) RegisterRemotes(remotes map[string]string) *CacheServer {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	remoteList := []string{}
	for identifier, addr := range remotes {
		// make sure server addr does not end with '/'
		if addr[len(addr)-1] == '/' {
			addr = addr[:len(addr)-1]
		}
		cs.remoteMap[identifier] = addr

		remoteList = append(remoteList, identifier)
	}

	cs.peerMgr.AddNodes(remoteList...)
	return cs
}
func (cs *CacheServer) RegisterRemote(identifier string, addr string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.remoteMap[identifier] = addr
	cs.peerMgr.AddNodes(identifier)
}

func (cs *CacheServer) GetNode(key string) (string, bool) {
	// Pick remote server according to key
	peer := cs.peerMgr.GetNode(key)
	if peer == "" {
		return "", false
	}

	peerAddr := cs.remoteMap[peer]
	segs := strings.Split(peerAddr, "//")
	if segs[len(segs)-1] == cs.addr {
		return "", true
	}

	return peerAddr, true
}
