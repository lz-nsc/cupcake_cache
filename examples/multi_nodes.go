package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/lz-nsc/cupcake_cache"
)

const defaultProxyPort = "8888"

// TODO: Haven't implement a external node management service
func startProxy(apiAddr string, group *cupcake_cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("start to run proxy at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil))
}

func main() {
	var port int
	var isProxy bool
	flag.IntVar(&port, "port", 8001, "Cupcake_cache node port")
	flag.BoolVar(&isProxy, "name", false, "name of this server")
	flag.BoolVar(&isProxy, "proxy", false, "Start a proxy or not")
	flag.Parse()
	var fakeDatabase = map[string]string{
		"Hello": "World",
		"Hi":    "Lynn",
		"Bye":   "History",
		"Good":  "Night",
		"Test":  "Done",
	}
	//TODO: Haven't implement dynamically regiest remote node yet
	addrMap := map[string]string{
		"node_1": "http://localhost:8080",
		"node_2": "http://localhost:8081",
		"node_3": "http://localhost:8082",
	}
	group := cupcake_cache.NewGroup("test", 0, cupcake_cache.GetterFunc(func(key string) ([]byte, error) {
		// get data from fake database
		if val, ok := fakeDatabase[key]; ok {
			return []byte(val), nil
		}
		return nil, nil
	}))

	if isProxy {
		proxyAddr := "localhost:" + defaultProxyPort
		go startProxy(proxyAddr, group)
	}

	remoteHandler := cupcake_cache.NewCacheHttp("localhost:"+strconv.Itoa(port), nil)

	group.SetRemote(remoteHandler.RegisterRemotes(addrMap))
	log.Fatal(remoteHandler.RunServer())
}
