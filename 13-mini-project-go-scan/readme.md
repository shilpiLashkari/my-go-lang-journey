# Day 13: Go-Scan - Concurrent Port Scanner

## Overview

A high-performance CLI port scanner that uses the **Worker Pool** pattern to scan thousands of ports in seconds. This project demonstrates advanced concurrency management and network programming in Go.

## Features

- **Worker Pool Architecture**: Efficiently manages a fixed number of goroutines to prevent resource exhaustion.
- **Port Range Scanning**: Specify custom ranges (e.g., `1-100`, `1-1024`).
- **Configurable Workers**: Adjust the number of concurrent workers for speed/stability.
- **Adjustable Timeout**: Set custom timeouts for network probes.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 13-mini-project-go-scan
   ```
2. Run with default settings (scans 1-1024 on scanme.nmap.org):
   ```bash
   go run main.go
   ```
3. Run with custom parameters:
   ```bash
   go run main.go -host 127.0.0.1 -ports 1-8080 -workers 500 -timeout 200ms
   ```

## Learning Reflection

- **Worker Pools**: Learned how to distribute tasks across a pool of goroutines using buffered channels.
- **TCP Networking**: Used `net.DialTimeout` to probe network ports.
- **Flag Parsing**: Leveraged the `flag` package for a professional CLI experience.
