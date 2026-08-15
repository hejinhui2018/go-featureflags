// Package ledgerbatch applies account balance changes atomically.
package ledgerbatch

import "errors"

// Entry changes one account balance by Delta units.
type Entry struct {
	Account string
	Delta   int
}

// ErrInsufficientFunds reports that a batch would make an account negative.
var ErrInsufficientFunds = errors.New("insufficient funds")

// Ledger stores account balances.
type Ledger struct {
	balances map[string]int
}

// New creates a ledger with a copy of initial balances.
func New(initial map[string]int) *Ledger {
	balances := make(map[string]int, len(initial))
	for account, balance := range initial {
		balances[account] = balance
	}
	return &Ledger{balances: balances}
}

// Snapshot returns a copy of the current balances.
func (l *Ledger) Snapshot() map[string]int {
	copyOfBalances := make(map[string]int, len(l.balances))
	for account, balance := range l.balances {
		copyOfBalances[account] = balance
	}
	return copyOfBalances
}

// ApplyBatch applies every entry or leaves the ledger unchanged. Entries for
// one account are evaluated in their input order, so their effects accumulate.
func (l *Ledger) ApplyBatch(entries []Entry) error {
	pending := make(map[string]int, len(l.balances))
	for account, balance := range l.balances {
		pending[account] = balance
	}

	for _, entry := range entries {
		next := pending[entry.Account] + entry.Delta
		if next < 0 {
			return ErrInsufficientFunds
		}
		pending[entry.Account] += entry.Delta
	}

	l.balances = pending
	return nil
}
