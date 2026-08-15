package main

import (
	"fmt"
	"os"

	"github.com/hejinhui2018/kvjournal-atomic-write-005/journal"
)

func main() {
	if len(os.Args) != 5 || os.Args[1] != "put" {
		fmt.Fprintln(os.Stderr, "usage: kvjournal put PATH KEY VALUE")
		os.Exit(2)
	}
	s, err := journal.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer s.Close()
	if err := s.Put(os.Args[3], os.Args[4]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
