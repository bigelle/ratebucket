package ratebucket

import (
	"sync"
)

const (
	defaultInitialTokens int64   = 1000
	defaultCapacity      int64   = 1000
	defaultRefillRate    float64 = 5
)

type PoolConfig struct {
	InitialTokens int64
	Capacity      int64
	RefillRate    float64
}

type Pool struct {
	mu      sync.Mutex
	buckets map[any]*Bucket
	config  PoolConfig
}

func NewPool() *Pool {
	return &Pool{
		buckets: make(map[any]*Bucket),
		config: PoolConfig{
			InitialTokens: defaultInitialTokens,
			Capacity:      defaultCapacity,
			RefillRate:    defaultRefillRate,
		},
	}
}

func NewPoolConfig(cfg PoolConfig) *Pool {
	return &Pool{
		buckets: make(map[any]*Bucket),
		config:  cfg,
	}
}

func (p *Pool) Allow(key any) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	b, ok := p.buckets[key]
	if !ok {
		b = NewBucket(
			WithInitialTokens(p.config.InitialTokens),
			WithCap(p.config.Capacity),
			WithRate(p.config.RefillRate),
		)
		p.buckets[key] = b
	}

	return b.Allow()
}
