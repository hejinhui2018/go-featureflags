package main

import (
	"fmt"

	"ledgerbatch"
)

func main() {
	ledger := ledgerbatch.New(map[string]int{"alice": 10})
	if err := ledger.ApplyBatch([]ledgerbatch.Entry{{Account: "alice", Delta: -3}}); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ledger.Snapshot())
}
