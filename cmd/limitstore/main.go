package main

import (
	"fmt"
	"os"

	"github.com/hejinhui2018/limitstore-capacity-update-007/limitstore"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: limitstore KEY VALUE")
		os.Exit(2)
	}
	s := limitstore.New(1)
	if err := s.Put(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
