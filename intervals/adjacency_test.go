package intervals

import (
	"reflect"
	"testing"
)

func TestMergeMergesTouchingIntervals(t *testing.T) {
	got := Merge([]Interval{{Start: 1, End: 3}, {Start: 3, End: 5}})
	want := []Interval{{Start: 1, End: 5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Merge() = %#v, want %#v", got, want)
	}
}
