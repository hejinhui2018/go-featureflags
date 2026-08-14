package main

import (
	"log"
	"net/http"

	"github.com/hejinhui2018/go-featureflags/internal/flags"
	"github.com/hejinhui2018/go-featureflags/internal/httpapi"
)

func main() {
	store := flags.NewMemoryStore(map[string]map[string]bool{
		"tenant-a": {"checkout": true, "search": true},
		"tenant-b": {"checkout": false, "search": true},
	})
	service := flags.NewService(store)
	handler := httpapi.NewHandler(service)

	log.Println("feature flag service listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
