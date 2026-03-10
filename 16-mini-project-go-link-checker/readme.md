# Day 16: Go-Link-Checker - Concurrent URL Verifier 🔗

## Overview

For Day 16, I've built **Go-Link-Checker**, a high-performance tool that verifies the status of multiple URLs concurrently. This project is a practical application of Go's concurrency strengths, allowing me to check hundreds of links in seconds. I've used worker pools to manage resources and the `net/http` package's `HEAD` method to speed up the verification process.

## Features

- **Concurrent Verification**: Uses a worker pool to scan multiple links at once without overwhelming my system or the network.
- **HEAD Requests**: Optimized for speed by fetching only response headers instead of full body content.
- **Statistics Dashboard**: Provides a summary of alive vs. broken links.
- **Configurable**:
  - Set custom number of workers for parallel processing.
  - Set network timeouts to skip unresponsive servers.
- **Input Flexibility**: Reads URLs from a simple text file.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 16-mini-project-go-link-checker
   ```
2. Create a `links.txt` file with one URL per line.
3. Run the checker:
   ```bash
   go run main.go -file links.txt -workers 20
   ```

## Learning Reflection

- **Worker Pools**: Solidified my understanding of how to use buffered channels to distribute work among a set of persistent goroutines.
- **HEAD vs. GET**: Learned that `http.Head` is significantly more efficient for link checking as it avoids downloading large pages.
- **Graceful Termination**: Used `sync.WaitGroup` to ensure all workers finish their tasks before the program exits and prints the summary.

---

_Keeping my links alive with the power of Go!_
