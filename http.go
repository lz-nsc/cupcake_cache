package cupcake_cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/lz-nsc/cupcake_cache/log"
)

const defaultBasePath = "/cupcakecache/"
const defaultReplicas = 4

type cacheHttp struct {
	addr      string
	basePath  string
	peerMgr   *Manager
	mu        sync.Mutex
	remoteMap map[string]string
}

func NewCacheHttp(addr string, hashFn HashFunc) *cacheHttp {
	return &cacheHttp{
		addr:      addr,
		basePath:  defaultBasePath,
		peerMgr:   NewManager(defaultReplicas, hashFn),
		remoteMap: map[string]string{},
	}
}

func (cs *cacheHttp) RunServer() error {
	http.Handle(cs.basePath, cs)
	return http.ListenAndServe(cs.addr, nil)
}

// cacheHttp accepte request /<basePath>/<groupName>/<key>
func (cs *cacheHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Infof("new request: %s %s", r.Method, r.URL.Path)

	routes := strings.Split(path[len(cs.basePath):], "/")
	if len(routes) != 2 {
		log.Infof("bad request, want 2 params, got: %d", len(routes))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	groupName := routes[0]
	key := routes[1]
	group := GetGroup(groupName)

	if group == nil {
		w.Write([]byte(fmt.Sprintf("gourp %s not found", groupName)))
		return
	}
	bv, err := group.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bv.ByteSlice())
}

func (cs *cacheHttp) RegisterRemotes(remotes map[string]string) *cacheHttp {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	remoteList := []string{}
	for identifier, addr := range remotes {
		// make sure server addr does not end with '/'
		if addr[len(addr)-1] == '/' {
			addr = addr[:len(addr)-1]
		}
		cs.remoteMap[identifier] = addr + cs.basePath

		remoteList = append(remoteList, identifier)
	}

	cs.peerMgr.AddNodes(remoteList...)
	return cs
}

func (cs *cacheHttp) RegisterRemote(identifier string, addr string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.remoteMap[identifier] = addr
	cs.peerMgr.AddNodes(identifier)
}

func (cs *cacheHttp) remoteGet(group string, key string) ([]byte, error) {
	// Pick remote server according to key
	peer := cs.peerMgr.GetNode(key)
	if peer == "" {
		return nil, fmt.Errorf("failed to find remote with peer key %s", key)
	}

	peerAddr := cs.remoteMap[peer]
	segs := strings.Split(peerAddr, "//")
	log.Debugf("segs: %v, addr: %s", segs, cs.addr)
	if segs[len(segs)-1] == cs.addr+cs.basePath {
		return nil, nil
	}

	log.Debugf("Successfully got remote peer, addr: %s", peerAddr)
	// Send request to remote server
	url := fmt.Sprintf(
		"%v%v/%v",
		peerAddr,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote peer returned: %v", resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return bytes, nil
}
