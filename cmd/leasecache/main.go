package main

import (
	"log"
	"net/http"
	"time"

	"leasecache"
)

func main() {
	store := leasecache.NewStore(nil)
	if err := store.Put("demo", []byte("active"), 5*time.Minute); err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", leasecache.Handler(store)))
}
