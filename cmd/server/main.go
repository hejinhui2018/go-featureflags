package main

import (
	"log"
	"net/http"

	"github.com/hejinhui2018/configstore-validation-004/internal/httpapi"
	"github.com/hejinhui2018/configstore-validation-004/internal/store"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", httpapi.NewHandler(store.NewMemory())))
}
