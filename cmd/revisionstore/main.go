package main

import (
	"fmt"

	"github.com/hejinhui2018/revisionstore-conflict-validation-010/revisionstore"
)

func main() {
	s := revisionstore.New()
	created, err := s.Put("theme", "dark", 0)
	if err != nil {
		panic(err)
	}
	updated, err := s.Put("theme", "light", created.Revision)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s revision=%d\n", updated.Value, updated.Revision)
}
