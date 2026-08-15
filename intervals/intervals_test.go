package intervals

import (
	"reflect"
	"testing"
)

func TestMergeNormalizesAndCoalesces(t *testing.T) {
	got := Merge([]Interval{{8, 10}, {1, 3}, {2, 4}, {12, 11}})
	want := []Interval{{1, 4}, {8, 10}, {11, 12}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Merge() = %#v, want %#v", got, want)
	}
}

func TestMergeEmpty(t *testing.T) {
	if got := Merge(nil); got != nil {
		t.Fatalf("Merge(nil) = %#v, want nil", got)
	}
}
