package main

import (
	"log"
	"net/http"

	"checkpointstore"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: checkpointstore.Handler(checkpointstore.NewStore()),
	}
	log.Fatal(server.ListenAndServe())
}
