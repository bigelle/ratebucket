package ratebucket

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBucketManager_Allow(t *testing.T) {
	p := NewPoolConfig(PoolConfig{
		Capacity:      3,
		InitialTokens: 3,
		RefillRate:    1,
	})

	key := "user123"

	require.True(t, p.Allow(key), "first request should pass")
	require.True(t, p.Allow(key), "second request should pass")
	require.True(t, p.Allow(key), "third request should pass")
	require.False(t, p.Allow(key), "fourth request should be rejected")

	time.Sleep(1100 * time.Millisecond)

	assert.True(t, p.Allow(key), "after refill one request should pass")
	assert.False(t, p.Allow(key), "next request should be rejected")
}

func TestBucketManager_AllowConcurrent(t *testing.T) {
	p := NewPoolConfig(PoolConfig{
		Capacity:      3,
		InitialTokens: 3,
		RefillRate:    1,
	})

	key := "concurrentUser"

	p.mu.Lock()
	p.buckets[key] = NewBucket(WithInitialTokens(5), WithRate(0))
	p.mu.Unlock()

	var wg sync.WaitGroup
	results := make(chan bool, 20)

	for range 20 {
		wg.Go(func() {
			results <- p.Allow(key)
		})
	}

	wg.Wait()
	close(results)

	var allowed, denied int
	for r := range results {
		if r {
			allowed++
		} else {
			denied++
		}
	}

	assert.Equal(t, 5, allowed, "exactly 5 requests should be allowed")
	assert.Equal(t, 15, denied, "remaining requests should be rejected")
}
