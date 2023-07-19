package main

import (
	"log"

	"github.com/lz-nsc/cupcake_cache"
)

var fakeDatabase = map[string]string{
	"Hello": "World",
	"Hi":    "Lynn",
	"Bye":   "History",
}

func main() {
	cupcake_cache.NewGroup("test", 0, cupcake_cache.GetterFunc(func(key string) ([]byte, error) {
		// get data from fake database
		if val, ok := fakeDatabase[key]; ok {
			return []byte(val), nil
		}
		return nil, nil
	}))

	server, err := cupcake_cache.NewServer("http", "", "localhost:8080", nil)
	if err != nil {
		panic("fail to create cache server")
	}
	log.Fatal(server.Run())
}
