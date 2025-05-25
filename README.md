# Go Concurrency & Thread-Safe LRU Cache
Ever tried building a web service only to find multiple users crashing your app due to bad thread handling? This project teaches you how to make concurrency work for you, not against you. Learn Go concurrency by building a thread-safe LRU cache from scratch. This project covers race conditions, synchronization, and performance optimization through hands-on coding.

# What You'll Build 
A production-ready LRU (Least Recently Used) cache that can handle thousands of concurrent operations safely. You'll start with broken code and fix it step by step.


---

## 🛠 Setup
```bash
git clone [text](https://github.com/kayceenuel/go_concurrency.git)
cd concurrency-project
go mod tidy


## 📂 Project Structure

concurrency-project/
├── exercises/         # Hands-on learning exercises
│   ├── buggy/        # Faulty concurrent code with race conditions ❌
│   ├── atomics/      # Fixed versions using atomic operations ✅
│   ├── mutex/        # Basic Mutex-based synchronization
│   ├── rwmutex/      # Optimized Read-Write Mutex handling

├── cache/            # Core LRU cache logic
│   ├── cache.go      # Implements thread-safe LRU eviction
│   ├── optimized.go  # RWMutex & sharding for better performance
│   ├── statistics.go # Tracks cache hit/miss metrics
│   ├── *_test.go     # Unit tests ensuring concurrency safety

├── examples/         # Demonstrations & benchmarks
│   ├── main.go       # Usage examples & performance benchmarking

├── docs/             # Documentation & diagrams
│   ├── concurrency_flowchart.png  # Visual explanation of cache workflow
│   ├── performance_comparison.md  # Benchmark results & analysis

├── README.md         # Project overview & instructions 📜
└── go.mod            # Go module dependencies

``` 
# Quick Start 
See the Problems 

```bash 
go run exercises/buggy/main.go
``` 
See the Solution
```bash 
go run exercises/atomics/main.go
 
``` 
Run full demo: 
```bash 
go run examples/main.go
```
Run all tests (including race detection):
```bash 
go test -race ./cache
``` 
## Core Features
- Thread-safe LRU Cache: Concurrent read/write operations without corruption

- Automatic eviction: Least recently used items removed when cache is full

### Multiple implementations:

🔹 Basic Mutex-based cache

🔹 Optimized RWMutex cache for read-heavy workloads

🔹 Sharded cache for high concurrency

- Detailed statistics tracking: Hits, misses, evictions, and access patterns

## Usage Example 
// Create a cache with max 100 items
cache := NewCache[string, int](100)

// Thread-safe operations
cache.Put("user:123", 42)
value, found := cache.Get("user:123")

// Performance statistics
stats := cache.GetStatistics()
fmt.Printf("Hit rate: %.2f%%\n", stats.GetHitRate() * 100)

## Learning Roadmap
### Race Conditions

- Why concurrent code breaks ❌

- How to spot race conditions 🔍

- Fixing with atomic operations & mutexes ✅

### Building the Cache

- Implement LRU logic ⚙️

- Make it thread-safe 🚀

- Track detailed statistics 📊

### Optimization

- Use RWMutex for better read performance ⚡

- Reduce lock contention with sharding 🏆

- Run benchmarks & analyze performance 📈

## Testing
Always test concurrent code with race detection:
```bash
# Basic tests
go test ./cache

# Race detection (catches bugs!)
go test -race ./cache

# Benchmarks
go test -bench=. ./cache
``` 
