# Day 20: Go-Load-Tester - Concurrent HTTP Benchmarker 🚀🔥

## Overview

For Day 20, I've built **Go-Load-Tester**, a powerful CLI tool designed to stress-test web servers. Since a major part of my Go journey focuses on performance and concurrency, I wanted to build a tool that actually measures speed! This project leverages Go's goroutines to fire hundreds of HTTP requests simultaneously at a target URL, acting like a simplified version of professional benchmarking tools built in Go, like `hey` or `vegeta`.

## Features

- **Massive Concurrency**: Spawns multiple "workers" to hit the server in parallel.
- **Configurable Load**: I can set exactly how many total requests to send and how many workers should run concurrently.
- **Detailed Metrics**: Tracks and displays:
  - Total time taken.
  - Number of successful responses (HTTP 2xx).
  - Number of failed requests.
  - Average response time per request.
- **Safe State Management**: Safely increments shared counters across dozens of threads using `sync/atomic` and mutexes.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 20-mini-project-go-load-tester
   ```
2. Run a load test against a public testing endpoint (e.g., 100 requests, 10 at a time):
   ```bash
   go run main.go -url https://httpbin.org/get -requests 100 -concurrency 10
   ```
*(Note: Be careful not to use this on servers you don't own to avoid causing a Denial of Service! 😅)*

## Learning Reflection

- **Worker Pools at Scale**: I used a combination of `sync.WaitGroup` and buffered channels to perfectly divide the workload among the requested number of concurrent workers.
- **Atomic Operations**: Instead of locking the whole app with Mutexes to count successes, I learned to use `atomic.AddInt64` for blisteringly fast, thread-safe counter increments.
- **Time Measurement**: Used `time.Since` to accurately profile both individual requests and the entire benchmark run.

---
*Testing the limits of the web, powered by Go!*
