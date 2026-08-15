package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hejinhui2018/cursorstore-validation-008/cursorstore"
)

func main() {
	s := cursorstore.New()
	for _, item := range []cursorstore.Entry{{Key: "alpha", Value: "1"}, {Key: "beta", Value: "2"}, {Key: "gamma", Value: "3"}} {
		s.Put(item.Key, item.Value)
	}
	entries, next, err := s.List("", 2)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = json.NewEncoder(os.Stdout).Encode(struct {
		Entries []cursorstore.Entry `json:"entries"`
		Next    string              `json:"next_cursor"`
	}{entries, next})
}
