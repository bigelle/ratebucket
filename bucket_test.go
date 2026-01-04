package ratebucket

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBucket_tokens(t *testing.T) {
	b := NewBucket(
		WithInitialTokens(5),
		WithCap(10),
		WithRate(2),
	)

	require.Equal(t, int64(5), b.tokens(), "initial tokens mismatch")

	b.mu.Lock()
	b.a_tokens.Store(4)
	b.mu.Unlock()

	time.Sleep(1100 * time.Millisecond)

	tokensBefore := b.tokens()
	assert.GreaterOrEqual(t, tokensBefore, int64(6), "tokens should refill after 1s")
	assert.LessOrEqual(t, tokensBefore, int64(10), "tokens should not exceed capacity")

	time.Sleep(2100 * time.Millisecond)

	tokensAfter := b.tokens()
	assert.Equal(t, int64(b.cap), tokensAfter, "tokens should cap at capacity")
}

func TestBucket_tokensConcurrent(t *testing.T) {
	b := NewBucket(
		WithCap(10),
		WithInitialTokens(5),
		WithRate(5),
	)

	var wg sync.WaitGroup
	results := make(chan int64, 100)

	for range 100 {
		wg.Go(func() {
			tokens := b.tokens()
			results <- tokens
		})
	}

	wg.Wait()
	close(results)

	for v := range results {
		require.GreaterOrEqual(t, v, int64(0))
		require.LessOrEqual(t, v, b.cap)
	}
}

func TestBucket_Allow(t *testing.T) {
	b := NewBucket(
		WithCap(3),
		WithInitialTokens(3),
		WithRate(1),
	)

	require.True(t, b.Allow(), "first request should pass")
	require.True(t, b.Allow(), "second request should pass")
	require.True(t, b.Allow(), "third request should pass")
	require.False(t, b.Allow(), "fourth request should be rejected")

	time.Sleep(1100 * time.Millisecond)
	assert.True(t, b.Allow(), "after refill one request should pass")
	assert.False(t, b.Allow(), "next request should still be rejected")
}

func TestBucket_AllowConcurrent(t *testing.T) {
	b := NewBucket(
		WithInitialTokens(5),
		WithCap(5),
		WithRate(0),
	)

	var wg sync.WaitGroup
	results := make(chan bool, 20)

	for range 20 {
		wg.Go(func() {
			ok := b.Allow()
			results <- ok
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
