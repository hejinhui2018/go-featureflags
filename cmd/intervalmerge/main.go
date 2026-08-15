package main

import (
	"encoding/json"
	"fmt"

	"example.com/intervalmerge/intervals"
)

func main() {
	input := []intervals.Interval{{Start: 5, End: 7}, {Start: 1, End: 3}, {Start: 3, End: 4}}
	encoded, err := json.Marshal(intervals.Merge(input))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(encoded))
}
