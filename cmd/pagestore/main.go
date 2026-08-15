package main

import (
	"fmt"

	"github.com/hejinhui2018/go-featureflags/pagestore/pagestore"
)

func main() {
	s := pagestore.New()
	s.Append("alpha")
	s.Append("beta")
	page, err := s.List(0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(page)
}
