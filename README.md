# Go Concurrency & Thread-Safe LRU Cache
Ever tried building a web service only to find multiple users crashing your app due to bad thread handling? This project teaches you how to make concurrency work for you, not against you. Learn Go concurrency by building a thread-safe LRU cache from scratch. This project covers race conditions, synchronization, and performance optimization through hands-on coding.

# What You'll Build 
A production-ready LRU (Least Recently Used) cache that can handle thousands of concurrent operations safely. You'll start with broken code and fix it step by step.


---

## 🛠 Setup
```bash
git clone <your-repo>
cd concurrency-project
go mod tidy


## 📂 Project Structure
```plaintext
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
