// Package idemqueue processes keyed requests exactly once.
package idemqueue

import (
	"context"
	"crypto/sha256"
	"errors"
	"sync"
)

var (
	// ErrMissingKey reports an empty idempotency key.
	ErrMissingKey = errors.New("missing idempotency key")
	// ErrKeyConflict reports that a key was reused for different input.
	ErrKeyConflict = errors.New("idempotency key reused with different payload")
)

// Result is the response saved for a successful request.
type Result struct {
	Status int
	Body   []byte
}

// Executor performs the operation associated with a request.
type Executor func(context.Context, []byte) (Result, error)

type record struct {
	digest [sha256.Size]byte
	result Result
}

// Processor stores successful results by idempotency key.
type Processor struct {
	mu      sync.Mutex
	records map[string]record
}

// NewProcessor creates an empty processor.
func NewProcessor() *Processor {
	return &Processor{records: make(map[string]record)}
}

// Process executes the first request for a key and replays its result for an
// identical payload. Reusing the key for different input returns a conflict.
func (p *Processor) Process(ctx context.Context, key string, payload []byte, execute Executor) (Result, error) {
	if key == "" {
		return Result{}, ErrMissingKey
	}

	digest := sha256.Sum256(payload)
	p.mu.Lock()
	defer p.mu.Unlock()

	if saved, ok := p.records[key]; ok {
		matchesPayload := saved.digest == digest
		if !matchesPayload {
			return Result{}, ErrKeyConflict
		}
		return cloneResult(saved.result), nil
	}

	result, err := execute(ctx, append([]byte(nil), payload...))
	if err != nil {
		return Result{}, err
	}
	result = cloneResult(result)
	p.records[key] = record{digest: digest, result: result}
	return cloneResult(result), nil
}

func cloneResult(result Result) Result {
	result.Body = append([]byte(nil), result.Body...)
	return result
}
