package ledgerbatch

import (
	"reflect"
	"testing"
)

func TestApplyBatchCommitsAllEntries(t *testing.T) {
	ledger := New(map[string]int{"alice": 10, "bob": 2})

	if err := ledger.ApplyBatch([]Entry{{Account: "alice", Delta: -3}, {Account: "bob", Delta: 4}}); err != nil {
		t.Fatalf("ApplyBatch() error = %v", err)
	}

	want := map[string]int{"alice": 7, "bob": 6}
	if got := ledger.Snapshot(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Snapshot() = %#v, want %#v", got, want)
	}
}

func TestApplyBatchIsAtomicOnError(t *testing.T) {
	ledger := New(map[string]int{"alice": 10})
	before := ledger.Snapshot()

	if err := ledger.ApplyBatch([]Entry{{Account: "alice", Delta: -6}, {Account: "alice", Delta: -5}}); err != ErrInsufficientFunds {
		t.Fatalf("ApplyBatch() error = %v, want %v", err, ErrInsufficientFunds)
	}
	if got := ledger.Snapshot(); !reflect.DeepEqual(got, before) {
		t.Fatalf("failed batch changed balances to %#v, want %#v", got, before)
	}
}

func TestApplyBatchUsesEntryOrder(t *testing.T) {
	ledger := New(map[string]int{"alice": 1})

	if err := ledger.ApplyBatch([]Entry{{Account: "alice", Delta: 4}, {Account: "alice", Delta: -3}}); err != nil {
		t.Fatalf("ApplyBatch() error = %v", err)
	}
	if got, want := ledger.Snapshot()["alice"], 2; got != want {
		t.Fatalf("alice balance = %d, want %d", got, want)
	}
}

func TestNewAndSnapshotDoNotShareInput(t *testing.T) {
	initial := map[string]int{"alice": 5}
	ledger := New(initial)
	initial["alice"] = 99

	if got, want := ledger.Snapshot()["alice"], 5; got != want {
		t.Fatalf("alice balance = %d, want %d", got, want)
	}
}
