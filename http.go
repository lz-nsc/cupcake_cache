package cupcake_cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/lz-nsc/cupcake_cache/log"
)

const defaultBasePath = "/cupcakecache/"

type cacheHttpServer struct {
	name        string
	addr        string
	basePath    string
	cacheServer *CacheServer
}

var _ Server = (*cacheHttpServer)(nil)

func init() {
	serverMap["http"] = NewCacheHttpServer
}

func NewCacheHttpServer(name string, addr string, cs *CacheServer) Server {
	log.WithServer(name, addr)
	return &cacheHttpServer{
		name:        name,
		addr:        addr,
		basePath:    defaultBasePath,
		cacheServer: cs,
	}
}

func (cs *cacheHttpServer) Run() error {
	http.Handle(cs.basePath, cs)
	return http.ListenAndServe(cs.addr, nil)
}

// cacheHttp accepte request /<basePath>/<groupName>/<key>
func (cs *cacheHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (cs *cacheHttpServer) RemoteGet(group string, key string) ([]byte, error) {
	// Pick remote server according to key

	peerAddr, ok := cs.cacheServer.GetNode(key)
	if !ok {
		return nil, fmt.Errorf("failed to get remote node with key %s", key)
	}
	if peerAddr == "" {
		return nil, nil
	}

	log.Debugf("Successfully got remote peer, addr: %s", peerAddr)
	// Send request to remote server
	url := fmt.Sprintf(
		"%v%s%v/%v",
		peerAddr,
		cs.basePath,
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
