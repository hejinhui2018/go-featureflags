package main

import (
	"fmt"
	"tokenbucket"
)

func main() {
	bucket := tokenbucket.New(4)
	_ = bucket.Reserve(3)
	bucket.Release(5)
	fmt.Println(bucket.Available())
}
