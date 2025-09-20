package ratebucket

import (
	"sync/atomic"
	"time"
)

type Bucket struct {
	a_tokens   atomic.Int64
	lastRefill time.Time
	cap        int64
	rate       float64
}

func NewBucket(opts ...BucketOption) *Bucket {
	buck := &Bucket{
		a_tokens:   atomic.Int64{},
		lastRefill: time.Now(),
		cap:        1000,
		rate:       5,
	}

	for _, opt := range opts {
		opt(buck)
	}

	return buck
}

type BucketOption func(*Bucket)

func WithInitialTokens(t int64) BucketOption {
	return func(b *Bucket) {
		b.a_tokens.Store(t)
	}
}

func WithRate(r float64) BucketOption {
	return func(b *Bucket) {
		b.rate = r
	}
}

func WithCap(c int64) BucketOption {
	return func(b *Bucket) {
		b.cap = c
	}
}

func (b *Bucket) tokens() int64 {
	now := time.Now()

	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed > 0 {
		refill := elapsed * b.rate
		b.a_tokens.Store(min(b.cap, b.a_tokens.Load()+int64(refill)))
		b.lastRefill = now
	}

	return b.a_tokens.Load()
}

func (b *Bucket) Allow() bool {
	t := b.tokens()
	if t-1 < 0 {
		return false
	}

	b.a_tokens.Store(t - 1)
	return true
}
