package main

import (
	"log"
	"net/http"

	"rangeserve"
)

func main() {
	http.Handle("/asset", rangeserve.Handler([]byte("example content\n")))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
