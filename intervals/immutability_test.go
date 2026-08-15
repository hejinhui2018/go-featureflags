package intervals

import (
	"reflect"
	"testing"
)

func TestMergeDoesNotModifyInput(t *testing.T) {
	input := []Interval{{Start: 8, End: 2}, {Start: 1, End: 4}}
	wantInput := append([]Interval(nil), input...)

	got := Merge(input)
	if want := []Interval{{Start: 1, End: 8}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Merge() = %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(input, wantInput) {
		t.Fatalf("Merge() changed input to %#v, want %#v", input, wantInput)
	}
}
