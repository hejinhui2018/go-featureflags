package main

import (
	"fmt"
	"os"

	"github.com/hejinhui2018/batchstore-atomic-import-006/batchstore"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: batchstore IMPORT.jsonl")
		os.Exit(2)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	if err := batchstore.New(nil).Import(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
