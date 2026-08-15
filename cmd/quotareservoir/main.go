package main

import (
	"log"
	"net/http"

	"quotareservoir"
)

func main() {
	store, err := quotareservoir.NewStore(100)
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Addr: ":8080", Handler: quotareservoir.Handler(store)}
	log.Fatal(server.ListenAndServe())
}
