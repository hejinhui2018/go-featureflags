package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"idemqueue"
)

func main() {
	var nextID int64
	processor := idemqueue.NewProcessor()
	http.Handle("/jobs", idemqueue.Handler(processor, func(_ context.Context, _ []byte) (idemqueue.Result, error) {
		id := atomic.AddInt64(&nextID, 1)
		return idemqueue.Result{Status: http.StatusCreated, Body: []byte(fmt.Sprintf("job-%d", id))}, nil
	}))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
