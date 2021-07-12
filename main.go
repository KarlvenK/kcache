package main

import (
	"fmt"
	"github.com/KarlvenK/kcache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "1",
	"Jack": "2",
	"Sam":  "3",
}

func main() {
	kcache.NewGroup("int", 2<<10, kcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := kcache.NewHTTPPool(addr)
	log.Println("kcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
