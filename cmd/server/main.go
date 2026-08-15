package main

import (
	"log"
	"net/http"

	"github.com/hejinhui2018/recordstore-validation-003/internal/httpapi"
	"github.com/hejinhui2018/recordstore-validation-003/internal/store"
)

func main() {
	memory := store.NewMemory(map[string]string{"welcome": "hello"})
	log.Fatal(http.ListenAndServe(":8080", httpapi.NewHandler(memory)))
}
