package main

import (
	"fmt"
	"time"

	"ttlcache"
)

func main() {
	now := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	cache := ttlcache.New(func() time.Time { return now })
	cache.Set("boundary", "stale", 0)
	value, ok := cache.Get("boundary")
	fmt.Printf("value=%q ok=%v\n", value, ok)
}
