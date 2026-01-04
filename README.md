# ⭐ Ratebucket

A small, fast, and idiomatic **Go rate limiting library** based on the token bucket algorithm.

**Ratebucket** is designed for developers who want precise rate control without heavyweight dependencies or complex configuration. It fits naturally into Go’s concurrency model and works equally well for services, workers, and CLI tools.

## 🚀 Key Features

- ⚡ **Token Bucket Algorithm** — Predictable and well-understood rate limiting semantics  
- 🧩 **Pool Support** — Efficient reuse of buckets for dynamic keys or high-cardinality workloads  
- 🪶 **Lightweight** — Minimal code, no external dependencies  
- 🧪 **Tested** — Includes unit tests for core behavior  
- 🧠 **Explicit API** — No hidden magic, easy to reason about

## 📦 Installation

```bash
go get github.com/bigelle/ratebucket
```

## 🔧 Usage Example (with Pool)

Pool is useful when you need many independent rate limiters, for example per user, per API key, or per resource.

```go
package main

import (
	"fmt"
	"time"

	"github.com/bigelle/ratebucket"
)

func main() {
	// Create a pool of buckets:
	//  - refill rate: 1 token per second
	//  - capacity: 5 tokens
	pool := ratebucket.NewPool(1, 5)

	key := "user:123"

	// Get a bucket for a specific key
	bucket := pool.Get(key)

	if bucket.Take() {
		fmt.Println("Request allowed")
	} else {
		fmt.Println("Rate limit exceeded")
	}

	// Buckets are automatically reused by the pool
	time.Sleep(time.Second)
}
```

This pattern is ideal for:
- per-user or per-client rate limiting
- API clients with multiple credentials
- background workers sharing common limits

## 🧠 How It Works

Ratebucket implements the classic token bucket algorithm:
- Each bucket has a fixed capacity.
- Tokens are replenished over time at a defined rate.
- An action is allowed only if a token can be consumed.

The Pool manages multiple buckets efficiently, letting you scale rate limiting across dynamic keys without manual bookkeeping.

## 🎯 Why Ratebucket?

Ratebucket is intentionally small.

It doesn’t try to be a framework — instead, it gives you reliable building blocks that you can compose into your own control logic. This makes it easy to audit, test, and adapt to different workloads.

If you need:
- deterministic behavior,
- low overhead,
- and full control over rate limiting logic,

Ratebucket is a solid choice.

## 📄 License

MIT License
