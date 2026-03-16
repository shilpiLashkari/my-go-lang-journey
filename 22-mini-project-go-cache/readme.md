# Day 22: Go-Cache - In-Memory Key-Value Store 🧠⚡

## Overview

For Day 22, I've built **Go-Cache**, a high-performance, thread-safe in-memory key-value store. Real-world applications rely heavily on caching systems (like Redis or Memcached) to speed up data retrieval. This project recreates the fundamental concepts behind those systems. I focused on Go's `sync` package and learned how to safely manage state across hundreds of concurrent goroutines while automatically cleaning up old data.

## Features

- **Blisteringly Fast**: Stores data directly in memory (RAM).
- **Thread-safe**: Uses `sync.RWMutex` to allow multiple goroutines to read simultaneously, while safely locking the map when data is being written or deleted.
- **Time-to-Live (TTL)**: Items can have an expiration time. If you wait too long, they disappear!
- **Background Janitor**: A background goroutine sweeps through the cache automatically at set intervals, actively deleting expired items to save memory.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 22-mini-project-go-cache
   ```
2. Run the simulation:
   ```bash
   go run main.go
   ```
3. Watch the console as concurrent workers perform rapid read/write operations, and observe the janitor cleaning up expired keys.

## Learning Reflection

- **Mutex vs. RWMutex**: I learned that `sync.Mutex` completely locks access for everyone, whereas `sync.RWMutex` allows multiple *readers* at the exact same time, massively boosting performance for read-heavy apps like caches.
- **Memory Leaks**: Discovered why a Background Janitor is critical. Without it, expired items would sit in the map forever (until requested), slowly eating up all server memory.
- **Graceful Ticker Usage**: Used `time.Ticker` inside a long-running goroutine to execute periodic cleanup tasks.

---
*Serving data faster than a database ever could!*
