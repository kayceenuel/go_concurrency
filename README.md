# Go Concurrency & Thread-Safe LRU Cache
Ever tried building a web service only to find multiple users crashing your app due to bad thread handling? This project teaches you how to make concurrency work for you, not against you. Learn Go concurrency by building a thread-safe LRU cache from scratch. This project covers race conditions, synchronization, and performance optimization through hands-on coding.

# What You'll Build 
A production-ready LRU (Least Recently Used) cache that can handle thousands of concurrent operations safely. You'll start with broken code and fix it step by step.

# Setup
git clone https://github.com/kayceenuel/go_concurrency.git
cd concurrency-project
go mod tidy

# Project Structure
concurrency-project/
├── exercises/
│   ├── buggy/main.go        # Broken concurrent code
│   └── atomics/main.go      # Fixed versions
├── cache/
│   ├── cache.go             # Main LRU implementation  
│   ├── statistics.go        # Thread-safe stats
│   └── *_test.go           # Tests
└── examples/main.go         # Usage examples