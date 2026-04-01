# Day 38: Go-Port-Scanner - High-Speed Network Auditor ⚡️🛡️

## Overview

For Day 38, I built **Go-Port-Scanner**, a high-speed, concurrent TCP port scanning utility. This project demonstrates the legendary concurrency model of Go (goroutines and channels) to build a security tool that can scan thousands of network ports in mere seconds—a task that would take minutes in traditional synchronous languages.

## Features

- **Blazing Fast Scanning**: Uses a worker-pool pattern to scan multiple ports simultaneously.
- **Configurable Range**: Easily scan a specific port range (e.g., 20-443) or even the full 1-65535 range.
- **Intelligent Hub**: A central results collector ensures that only open ports are tracked and reported in a clean, sorted list.
- **Robust Connection Handling**: Implements `net.DialTimeout` to prevent hanging on "filtered" or non-responsive ports.
- **Pure Standard Library**: High-level network auditing built with zero external dependencies.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 38-mini-project-port-scanner
   ```
2. **Scan a target (e.g., scanme.nmap.org or localhost)**:
   ```bash
   # Scan ports 1 to 1024 on a public test server
   go run main.go -host scanme.nmap.org -ports 1-1024 -concurrency 100
   ```
3. Watch the real-time results in your terminal!

## Learning Reflection

- **Fan-out/Fan-in Pattern**: Implemented a classic concurrency architectural pattern where tasks are fanned out to workers and results are fanned in to a single collector.
- **Networking Primitives**: Gained a deeper understanding of TCP connection handshakes and timeout management using `net.DialTimeout`.
- **Systematic Resource Management**: Used `sync.WaitGroup` to ensure the program only exits after every worker has finished its task.

---
*Scanning for security, one port at a time. 🐹🛡️*
