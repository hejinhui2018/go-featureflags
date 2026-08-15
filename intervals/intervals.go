// Package intervals provides operations on closed integer intervals.
package intervals

// Interval is a closed interval [Start, End].
type Interval struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Merge returns normalized, sorted intervals. Overlapping and touching
// intervals are coalesced into one closed interval.
func Merge(input []Interval) []Interval {
	if len(input) == 0 {
		return nil
	}

	work := append([]Interval(nil), input...)
	for i := range work {
		if work[i].Start > work[i].End {
			work[i].Start, work[i].End = work[i].End, work[i].Start
		}
	}
	for i := 1; i < len(work); i++ {
		for j := i; j > 0 && work[j].Start < work[j-1].Start; j-- {
			work[j], work[j-1] = work[j-1], work[j]
		}
	}

	merged := make([]Interval, 0, len(work))
	current := work[0]
	for _, next := range work[1:] {
		if next.Start <= current.End {
			if next.End > current.End {
				current.End = next.End
			}
			continue
		}
		merged = append(merged, current)
		current = next
	}
	return append(merged, current)
}
