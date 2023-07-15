package cupcake_cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/cupcakecache/"

type cacheServer struct {
	addr     string
	basePath string
}

func NewCacheServer(addr string) *cacheServer {
	return &cacheServer{
		addr:     addr,
		basePath: defaultBasePath,
	}
}

func (cs *cacheServer) Run() error {
	http.Handle(cs.basePath, cs)
	return http.ListenAndServe(cs.addr, nil)
}

// cacheServer accepte request /<basePath>/<groupName>/<key>
func (cs *cacheServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("new request: %s %s", r.Method, r.URL.Path)

	routes := strings.Split(path[len(cs.basePath):], "/")
	if len(routes) != 2 {
		log.Printf("bad request, want 2 params, got: %d", len(routes))
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
