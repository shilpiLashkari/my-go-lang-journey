# Day 37: Go-HTTP-Bench - High Performance Load Tester ⚡️📈

## Overview

For Day 37, I built **Go-HTTP-Bench**, a lightweight and remarkably fast HTTP benchmarking utility. Go is the language behind modern cloud infrastructure (Docker, Kubernetes, Terraform) primarily because of its stellar performance. This tool is a pure-Go alternative to `ab` or `wrk`, designed to stress-test your web services with concurrent requests and provide clear latency reporting.

## Features

- **Concurrent Execution**: Spawns multiple "worker" goroutines to hammer your target URL simultaneously.
- **Detailed Statistics**: Reports Requests per Second (RPS), Total Requests, Success/Failure counts, and a breakdown of latency (Max, Min, Avg).
- **Graceful Control**: Limits the exact number of requests to prevent overwhelming the target beyond your intent.
- **Status Reporting**: Tracks and displays a breakdown of HTTP status codes (2xx, 4xx, 5xx).
- **Clean CLI UI**: Uses formatted strings to produce an easy-to-read performance summary.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 37-mini-project-http-bench
   ```
2. **Benchmark a local or remote URL**:
   ```bash
   # 1000 requests, 10 at a time
   go run main.go -url https://go.dev -requests 1000 -concurrency 10
   ```
3. Read the output report for the results!

## Learning Reflection

- **WaitGroups & Mutexes**: Practiced synchronizing workers using `sync.WaitGroup` and protecting shared statistical counters with `sync.Mutex`.
- **Measuring Latency**: Gained proficiency in using high-resolution timing (`time.Since`) to measure the performance of network calls.
- **Worker Pool Pattern**: Implemented a classic Go concurrency pattern where a pool of long-running workers processes tasks from a shared channel.

---
*Speed through concurrency. 🐹⚡️*
