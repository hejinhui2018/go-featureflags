package tokenbucket

import "errors"

var ErrInsufficientTokens = errors.New("insufficient tokens")

type Bucket struct {
	capacity int
	used     int
}

func New(capacity int) *Bucket {
	return &Bucket{capacity: capacity}
}

func (b *Bucket) Available() int {
	return b.capacity - b.used
}

func (b *Bucket) Reserve(n int) error {
	if n <= 0 {
		return nil
	}
	if b.used+n > b.capacity {
		return ErrInsufficientTokens
	}
	b.used += n
	return nil
}

func (b *Bucket) Release(n int) {
	if n <= 0 {
		return
	}
	b.used -= n
}