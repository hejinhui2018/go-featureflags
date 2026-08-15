package main

import (
	"fmt"

	"github.com/hejinhui2018/ttlstore-duration-validation-009/ttlstore"
)

func main() {
	s := ttlstore.New(nil)
	if err := s.Put("theme", "dark", 0); err != nil {
		panic(err)
	}
	value, err := s.Get("theme")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
