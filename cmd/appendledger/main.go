package main

import (
	"log"
	"net/http"

	"appendledger"
)

func main() {
	server := &http.Server{Addr: ":8080", Handler: appendledger.Handler(appendledger.NewStore())}
	log.Fatal(server.ListenAndServe())
}
